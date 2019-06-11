package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/stretchr/testify/require"
)

func TestConfigurePinCode(t *testing.T) {
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
		contentType                   string
		httpBody                      string
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
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  ContentTypeForm,
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodPost,
			status: http.StatusConflict,
			httpBody: toJSON(t, &ConfigurePinCodeRequest{
				RemovePin: false,
			}),
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
			httpBody: toJSON(t, &ConfigurePinCodeRequest{
				RemovePin: true,
			}),
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK remove pin false",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: toJSON(t, &ConfigurePinCodeRequest{
				RemovePin: false,
			}),
			gatewayConfigurePinCodeResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "configure pin code success msg",
			},
		},

		{
			name:   "200 - OK remove pin true",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: toJSON(t, &ConfigurePinCodeRequest{
				RemovePin: true,
			}),
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

			var body ConfigurePinCodeRequest
			err := json.Unmarshal([]byte(tc.httpBody), &body)
			if err == nil {
				gateway.On("ChangePin", newBoolPtr(body.RemovePin)).Return(tc.gatewayConfigurePinCodeResult, nil)
			}

			req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Set("Content-Type", contentType)

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
