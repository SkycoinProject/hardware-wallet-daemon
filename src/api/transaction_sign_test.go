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

func TestSignTransaction(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                         string
		method                       string
		status                       int
		contentType                  string
		httpBody                     string
		gatewaySignTransactionResult wire.Message
		httpResponse                 HTTPResponse
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
			name:        "400 - Missing Inputs",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "inputs is required"),
		},

		{
			name:        "400 - Missing InputIndexes",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "input_indexes is required"),
		},

		{
			name:        "400 - Missing AddressIndexes",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2", "3"},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "address_indexes is required"),
		},

		{
			name:        "400 - Missing Coins",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "coins is required"),
		},

		{
			name:        "400 - Missing Hours",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours is required"),
		},

		{
			name:        "400 - Missing Output Addresses",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusBadRequest,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:   []uint32{0, 1},
				Coins:          []string{"2", "3"},
				Hours:          []string{"2", "3"},
				AddressIndexes: []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "output_addresses is required"),
		},

		{
			name:        "422 - Input - InputIndexes mismatch",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "inputs length not equal to input_indexes length"),
		},

		{
			name:        "422 - OutputAddresses - Coins mismatch",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "output_addresses length not equal to coins length"),
		},

		{
			name:        "422 - OutputAddresses - Hours mismatch",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusUnprocessableEntity,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2"},
				AddressIndexes:  []int{0, 1},
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "output_addresses length not equal to hours length"),
		},

		{
			name:        "409 - Failure msg",
			method:      http.MethodPost,
			contentType: ContentTypeJSON,
			status:      http.StatusConflict,
			httpBody: toJSON(t, &TransactionSignRequest{
				Inputs: []string{
					"c2244e4912330d201d979f80db4df42118e49704e500e2e00a52a61954e8c663",
					"4f7250b0b1f588c4dedd5a4be984fab7215a773773480d8698e8f5ff04ef2611"},
				InputIndexes:    []uint32{0, 1},
				OutputAddresses: []string{"2M9hQ4LqEsBF5JZ3uBatnkaMgg9pN965JvG", "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"},
				Coins:           []string{"2", "3"},
				Hours:           []string{"2", "3"},
				AddressIndexes:  []int{0, 1},
			}),
			gatewaySignTransactionResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/transactionSign"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				var body TransactionSignRequest
				err := json.Unmarshal([]byte(tc.httpBody), &body)
				if err == nil {
					ins, outs, err := body.TransactionParams()
					if err == nil {
						gateway.On("TransactionSign", ins, outs).Return(tc.gatewaySignTransactionResult, nil)
					}
				}

				req, err := http.NewRequest(tc.method, "/api"+endpoint, strings.NewReader(tc.httpBody))
				require.NoError(t, err)

				contentType := tc.contentType
				if contentType == "" {
					contentType = ContentTypeJSON
				}

				req.Header.Set("Content-Type", contentType)

				rr := httptest.NewRecorder()
				handler := newServerMux(gateway, gateway)
				handler.ServeHTTP(rr, req)

				status := rr.Code
				require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

				var rsp HTTPResponse
				err = json.NewDecoder(rr.Body).Decode(&rsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Error, rsp.Error)

				if rsp.Data == nil {
					require.Nil(t, tc.httpResponse.Data)
				} else {
					require.NotNil(t, tc.httpResponse.Data)
				}
			})
		}
	}
}
