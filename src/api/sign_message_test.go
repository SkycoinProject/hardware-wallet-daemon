package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	"github.com/stretchr/testify/require"
)

func TestSignMessage(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                     string
		method                   string
		status                   int
		contentType              string
		httpBody                 string
		gatewaySignMessageResult wire.Message
		httpResponse             HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  ContentTypeForm,
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:        "422 - AddressN negative",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &SignMessageRequest{
				AddressN: -1,
				Message:  "foo",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be negative"),
		},

		{
			name:         "400 - empty message",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpBody:     toJSON(t, &SignMessageRequest{}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "message is required"),
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			status:      http.StatusConflict,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, &SignMessageRequest{
				AddressN: 0,
				Message:  "foo",
			}),
			gatewaySignMessageResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/sign_message"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				var body SignMessageRequest
				err := json.Unmarshal([]byte(tc.httpBody), &body)
				if err == nil {
					gateway.On("SignMessage", body.AddressN, body.Message).Return(tc.gatewaySignMessageResult, nil)
				}

				req, err := http.NewRequest(tc.method, "/api/v1"+endpoint, strings.NewReader(tc.httpBody))
				require.NoError(t, err)

				contentType := tc.contentType
				if contentType == "" {
					contentType = ContentTypeJSON
				}

				req.Header.Set("Content-Type", contentType)

				rr := httptest.NewRecorder()
				handler := newServerMux(defaultMuxConfig(), gateway, gateway)
				handler.ServeHTTP(rr, req)

				status := rr.Code
				require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

				var rsp HTTPResponse
				err = json.NewDecoder(rr.Body).Decode(&rsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Error, rsp.Error)

				if rsp.Data == nil {
					require.Nil(t, tc.httpResponse.Data)
				}
			})
		}
	}
}
