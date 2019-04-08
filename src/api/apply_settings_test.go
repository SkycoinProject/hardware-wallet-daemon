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

func TestApplySettings(t *testing.T) {
	type httpBody struct {
		usePassphrase string
		label         string
	}

	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                       string
		method                     string
		status                     int
		contentType                string
		httpBody                   *httpBody
		label                      string
		usePassphrase              bool
		gatewayApplySettingsResult wire.Message
		httpResponse               HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "400 - invalid passphrase",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for use-passphrase"),
			httpBody: &httpBody{
				usePassphrase: "foo",
			},
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusConflict,
			httpBody: &httpBody{
				label:         "foo",
				usePassphrase: "false",
			},
			label:         "foo",
			usePassphrase: false,
			gatewayApplySettingsResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: &httpBody{
				label:         "foo",
				usePassphrase: "false",
			},
			label:         "foo",
			usePassphrase: false,
			gatewayApplySettingsResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "success msg",
			},
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/applySettings"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				v := url.Values{}
				if tc.httpBody != nil {
					if tc.httpBody.usePassphrase != "" {
						v.Add("use-passphrase", tc.httpBody.usePassphrase)
					}

					if tc.httpBody.label != "" {
						v.Add("label", tc.httpBody.label)
					}
					if len(v) > 0 {
						endpoint += "?" + v.Encode()
					}
				}

				gateway.On("ApplySettings", tc.usePassphrase, tc.label).Return(tc.gatewayApplySettingsResult, nil)

				req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				handler := newServerMux(defaultMuxConfig(), gateway, gateway)
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
