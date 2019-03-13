package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	"github.com/stretchr/testify/require"
)

func TestFeatures(t *testing.T) {
	failureMsg := messages.Failure{
		Code:    messages.FailureType_Failure_NotInitialized.Enum(),
		Message: newStrPtr("failure msg"),
	}

	failureMsgBytes, err := failureMsg.Marshal()
	require.NoError(t, err)

	var featuresMsg = &messages.Features{
		Vendor: newStrPtr("Skycoin Foundation"),
	}

	featuresMsgBytes, err := featuresMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                  string
		method                string
		status                int
		contentType           string
		gatewayFeaturesResult wire.Message
		httpResponse          HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodPost,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:   "409 - Failure msg",
			method: http.MethodGet,
			status: http.StatusConflict,
			gatewayFeaturesResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: failureMsgBytes,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusConflict, "failure msg"),
		},

		{
			name:   "200 - OK",
			method: http.MethodGet,
			status: http.StatusOK,
			gatewayFeaturesResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Features),
				Data: featuresMsgBytes,
			},
			httpResponse: HTTPResponse{
				Data: FeaturesResponse{
					Features: featuresMsg,
				},
			},
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/features"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				gateway.On("GetFeatures").Return(tc.gatewayFeaturesResult, nil)

				req, err := http.NewRequest(tc.method, "/api"+endpoint, nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				handler := newServerMux(gateway, gateway)
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
					var resp FeaturesResponse
					err = json.Unmarshal(rsp.Data, &resp)
					require.NoError(t, err)

					require.Equal(t, tc.httpResponse.Data.(FeaturesResponse), resp)
				}
			})
		}
	}
}
