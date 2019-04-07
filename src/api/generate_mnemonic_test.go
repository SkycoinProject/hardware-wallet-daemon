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

func TestGenerateMnemonic(t *testing.T) {
	type httpBody struct {
		usePassphrase string
		wordCount     string
	}

	successMsg := messages.Success{
		Message: newStrPtr("Mnemonic successfully configured"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                          string
		method                        string
		status                        int
		httpBody                      *httpBody
		usePassphrase                 bool
		wordCount                     uint32
		httpResponse                  HTTPResponse
		gatewayGenerateMnemonicResult wire.Message
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
				wordCount:     "12",
			},
		},

		{
			name:         "400 - missing word-count",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing word-count"),
			httpBody: &httpBody{
				usePassphrase: "true",
			},
		},

		{
			name:         "400 - invalid word-count",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid value foo for word-count"),
			httpBody: &httpBody{
				usePassphrase: "true",
				wordCount:     "foo",
			},
		},

		{
			name:         "409 - Failure msg",
			method:       http.MethodPost,
			status:       http.StatusConflict,
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
			httpBody: &httpBody{
				wordCount:     "12",
				usePassphrase: "false",
			},
			gatewayGenerateMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			wordCount:     12,
			usePassphrase: false,
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			httpResponse: HTTPResponse{
				Data: *successMsg.Message,
			},
			httpBody: &httpBody{
				wordCount:     "12",
				usePassphrase: "false",
			},
			gatewayGenerateMnemonicResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			wordCount:     12,
			usePassphrase: false,
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				gateway := &MockGatewayer{}
				endpoint := "/generateMnemonic"

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				v := url.Values{}
				if tc.httpBody != nil {
					if tc.httpBody.usePassphrase != "" {
						v.Add("use-passphrase", tc.httpBody.usePassphrase)
					}

					if tc.httpBody.wordCount != "" {
						v.Add("word-count", tc.httpBody.wordCount)
					}
					if len(v) > 0 {
						endpoint += "?" + v.Encode()
					}
				}

				gateway.On("GenerateMnemonic", tc.wordCount, tc.usePassphrase).Return(tc.gatewayGenerateMnemonicResult, nil)

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
