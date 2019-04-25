package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/stretchr/testify/require"
)

func TestConfigurePinCode(t *testing.T) {
	type httpBody struct {
		removePin string
	}

	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("configure pin code success msg"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                          string
		method                        string
		status                        int
		httpBody                      *httpBody
		removePin                     bool
		gatewayConfigurePinCodeResult wire.Message
		httpResponse                  HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "400 - invalid remove_pin",
			method:       http.MethodPost,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for remove_pin"),
			httpBody: &httpBody{
				removePin: "foo",
			},
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodPost,
			status: http.StatusConflict,
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "409 - Failure msg with remove pin",
			method: http.MethodPost,
			status: http.StatusConflict,
			httpBody: &httpBody{
				removePin: "true",
			},
			removePin: true,
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodPost,
			status: http.StatusOK,
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "configure pin code success msg",
			},
		},

		{
			name:   "200 - OK with remove pin",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: &httpBody{
				removePin: "true",
			},
			removePin: true,
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "configure pin code success msg",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/configure_pin_code"
			gateway := &MockGatewayer{}

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.removePin != "" {
					v.Add("remove_pin", tc.httpBody.removePin)
				}

				if len(v) > 0 {
					endpoint += "?" + v.Encode()
				}
			}

			gateway.On("ChangePin", newBoolPtr(tc.removePin)).Return(tc.gatewayConfigurePinCodeResult, nil)

			req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
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
