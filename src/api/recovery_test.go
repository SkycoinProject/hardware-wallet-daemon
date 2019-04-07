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

func TestRecovery(t *testing.T) {
	type httpBody struct {
		usePassphrase string
		wordCount     string
		dryRun        string
	}

	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("recovery success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                  string
		method                string
		status                int
		httpBody              *httpBody
		usePassphrase         bool
		wordCount             uint32
		dryRun                bool
		httpResponse          HTTPResponse
		gatewayRecoveryResult wire.Message
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
				wordCount:     "2",
				dryRun:        "true",
			},
		},

		{
			name:         "400 - missing word-count",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing word-count"),
			httpBody: &httpBody{
				usePassphrase: "true",
				dryRun:        "true",
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
				dryRun:        "true",
			},
		},

		{
			name:         "400 - invalid dry-run",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for dry-run"),
			httpBody: &httpBody{
				usePassphrase: "false",
				wordCount:     "2",
				dryRun:        "foo",
			},
			wordCount: 2,
		},

		{
			name:         "409 - Failure msg",
			method:       http.MethodPost,
			status:       http.StatusConflict,
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
			httpBody: &httpBody{
				usePassphrase: "false",
				wordCount:     "2",
				dryRun:        "true",
			},
			gatewayRecoveryResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			wordCount:     2,
			dryRun:        true,
			usePassphrase: false,
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			httpResponse: HTTPResponse{
				Data: "recovery success msg",
			},
			httpBody: &httpBody{
				usePassphrase: "false",
				wordCount:     "2",
				dryRun:        "true",
			},
			gatewayRecoveryResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			wordCount:     2,
			dryRun:        true,
			usePassphrase: false,
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/recovery"
				gateway := &MockGatewayer{}

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

					if tc.httpBody.dryRun != "" {
						v.Add("dry-run", tc.httpBody.dryRun)
					}

					if len(v) > 0 {
						endpoint += "?" + v.Encode()
					}
				}

				gateway.On("Recovery", tc.wordCount, tc.usePassphrase, tc.dryRun).Return(tc.gatewayRecoveryResult, nil)

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
