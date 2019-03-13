package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	"github.com/stretchr/testify/require"
)

func TestSetMnemonic(t *testing.T) {
	type httpBody struct {
		mnemonic string
	}

	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("setmnemonic success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                     string
		method                   string
		status                   int
		contentType              string
		httpBody                 *httpBody
		gatewaySetMnemonicResult wire.Message
		httpResponse             HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "400 - missing mnemonic",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "mnemonic is required"),
			httpBody:     &httpBody{},
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodPost,
			status: http.StatusConflict,
			gatewaySetMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpBody: &httpBody{
				mnemonic: "cloud flower upset remain green metal below cup stem infant art thank",
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			gatewaySetMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpBody: &httpBody{
				mnemonic: "cloud flower upset remain green metal below cup stem infant art thank",
			},
			httpResponse: HTTPResponse{
				Data: "setmnemonic success msg",
			},
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/setMnemonic"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				v := url.Values{}
				if tc.httpBody != nil {
					if tc.httpBody.mnemonic != "" {
						v.Add("mnemonic", tc.httpBody.mnemonic)
					}

					if len(v) > 0 {
						endpoint += "?" + v.Encode()
					}
				}

				if tc.httpBody != nil {
					gateway.On("SetMnemonic", tc.httpBody.mnemonic).Return(tc.gatewaySetMnemonicResult, nil)
				}

				req, err := http.NewRequest(tc.method, "/api"+endpoint, nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				handler := newServerMux(gateway, gateway)
				handler.ServeHTTP(rr, req)

				status := rr.Code
				require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

				var rsp ReceivedHTTPResponse
				err = json.NewDecoder(rr.Body).Decode(&rsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Error, rsp.Error)

				if rsp.Data == nil {
					require.Nil(t, tc.httpResponse.Data)
				} else {
					require.NotNil(t, tc.httpResponse.Data)

					var resp string
					err = json.Unmarshal(rsp.Data, &resp)
					require.NoError(t, err)

					require.Equal(t, tc.httpResponse.Data.(string), resp)
				}
			})
		}
	}
}
