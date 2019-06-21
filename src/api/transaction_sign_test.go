package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
	"github.com/stretchr/testify/require"
)

func TestSignTransaction(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	successMsg := messages.Success{
		Message: newStrPtr("transaction sign success"),
	}

	successMsgBytes, err := successMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                         string
		method                       string
		status                       int
		contentType                  string
		httpBody                     string
		gatewaySignTransactionResult wire.Message
		err                          string
		httpResponse                 HTTPResponse
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
			name:        "400 - Input Hash Empty",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), ""}, {newUint32Ptr(1), ""},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "2", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			err:          "input hash cannot be empty",
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "input hash cannot be empty"),
		},

		{
			name:        "400 - Missing Coins",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Hours: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "coins cannot be empty"),
		},

		{
			name:        "400 - Missing Hours",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours cannot be empty"),
		},

		{
			name:        "400 - Missing Output Addresses",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Coins: "2", Hours: "2"},
					{Coins: "3", Hours: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "address cannot be empty"),
		},

		{
			name:        "422 - Invalid checksum",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsas5JZ3uBatnkaMgg9pN965JvG", Coins: "2", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, cipher.ErrAddressInvalidChecksum.Error()),
		},

		{
			name:        "422 - Missing Output Addresses",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "0.000000001010111001", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, droplet.ErrTooManyDecimals.Error()),
		},

		{
			name:        "422 - Missing Output Addresses",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "1", Hours: "0.2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "strconv.ParseUint: parsing \"0.2\": invalid syntax"),
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusConflict,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{newUint32Ptr(0), "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{newUint32Ptr(1), "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "2", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			gatewaySignTransactionResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			err:          "failure msg",
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:        "200 - Input Index Empty",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusOK,
			httpBody: toJSON(t, &TransactionSignRequest{
				TransactionInputs: []TransactionInput{
					{nil, "c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663"},
					{nil, "4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				},
				TransactionOutputs: []TransactionOutput{
					{Address: "2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", Coins: "2", Hours: "2"},
					{Address: "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8", Coins: "3", Hours: "3"},
				},
			}),
			gatewaySignTransactionResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Success),
				Data: successMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "transaction sign success",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/transaction_sign"
			gateway := &MockGatewayer{}

			if tc.httpBody != "" {
				var body TransactionSignRequest
				err := json.Unmarshal([]byte(tc.httpBody), &body)
				if err == nil {
					ins, outs, err := body.TransactionParams()
					if err == nil {
						gateway.On("TransactionSign", ins, outs).Return(tc.gatewaySignTransactionResult, nil)
					}
				}
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

			var rsp HTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
				if tc.err != "" {
					require.Equal(t, tc.err, rsp.Error.Message)
				}
			} else {
				require.NotNil(t, tc.httpResponse.Data)
			}
		})
	}
}
