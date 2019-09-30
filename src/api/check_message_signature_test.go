package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SkycoinProject/hardware-wallet-go/src/skywallet/wire"
	messages "github.com/SkycoinProject/hardware-wallet-protob/go"
	"github.com/stretchr/testify/require"
)

func TestCheckMessageSignature(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                               string
		method                             string
		status                             int
		contentType                        string
		httpBody                           string
		gatewayCheckMessageSignatureResult wire.Message
		httpResponse                       HTTPResponse
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
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:        "400 - Address missing",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Message:   "foo",
				Signature: "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "address is required"),
		},

		{
			name:        "422 - Address invalid",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Address:   "ca",
				Message:   "foo",
				Signature: "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "Invalid address length"),
		},

		{
			name:        "422 - Address invalid",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Address:   "c0",
				Message:   "foo",
				Signature: "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "Invalid base58 character"),
		},

		{
			name:        "400 - Signature missing",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Address: "u37EnnuQ4g58sWpd5Ns3FWGPwSgEuQGFBd",
				Message: "foo",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "signature is required"),
		},

		{
			name:        "400 - Message missing",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Address:   "u37EnnuQ4g58sWpd5Ns3FWGPwSgEuQGFBd",
				Signature: "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "message is required"),
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			status:      http.StatusConflict,
			contentType: ContentTypeJSON,
			httpBody: toJSON(t, &CheckMessageSignatureRequest{
				Address:   "u37EnnuQ4g58sWpd5Ns3FWGPwSgEuQGFBd",
				Signature: "GvKS4S3CA2YTpEPFA47yFdC5CP3y3qB18jwiX1URXqWQTvMjokd3A4upPz4wyeAyKJEtRdRDGUvUgoGASpsTTUeMn",
				Message:   "Hello World!",
			}),
			gatewayCheckMessageSignatureResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/check_message_signature"
			gateway := &MockGatewayer{}

			var body CheckMessageSignatureRequest
			err := json.Unmarshal([]byte(tc.httpBody), &body)
			if err == nil {
				gateway.On("CheckMessageSignature", body.Message, body.Signature, body.Address).Return(tc.gatewayCheckMessageSignatureResult, nil)
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
