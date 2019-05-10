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

func TestGenerateAddresses(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	responseAddressMsg := messages.ResponseSkycoinAddress{
		Addresses: []string{"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs"},
	}

	responseMsgBytes, err := responseAddressMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                    string
		method                  string
		status                  int
		contentType             string
		httpBody                string
		gatewayAddressGenResult wire.Message
		httpResponse            HTTPResponse
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
			name:         "422 - AddressN 0",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusUnprocessableEntity,
			httpBody:     toJSON(t, &GenerateAddressesRequest{}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be 0"),
		},

		{
			name:        "422 - AddressN negative",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &GenerateAddressesRequest{
				AddressN: -2,
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be negative"),
		},

		{
			name:        "422 - StartIndex negative",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &GenerateAddressesRequest{
				AddressN:   2,
				StartIndex: -2,
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "start_index cannot be negative"),
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusConflict,
			httpBody: toJSON(t, &GenerateAddressesRequest{
				AddressN:   2,
				StartIndex: 0,
			}),
			gatewayAddressGenResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:        "200 - OK",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusOK,
			httpBody: toJSON(t, &GenerateAddressesRequest{
				AddressN:   2,
				StartIndex: 0,
			}),
			gatewayAddressGenResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_ResponseSkycoinAddress),
				Data: responseMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: []string{"2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw", "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/generate_addresses"
			gateway := &MockGatewayer{}

			var body GenerateAddressesRequest
			err := json.Unmarshal([]byte(tc.httpBody), &body)
			if err == nil {
				gateway.On("AddressGen", body.AddressN, body.StartIndex, body.ConfirmAddress).Return(tc.gatewayAddressGenResult, nil)
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
				var resp []string
				err = json.Unmarshal(rsp.Data, &resp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.([]string), resp)
			}
		})
	}
}

func toJSON(t *testing.T, r interface{}) string {
	b, err := json.Marshal(r)
	require.NoError(t, err)
	return string(b)
}
