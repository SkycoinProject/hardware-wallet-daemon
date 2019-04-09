package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	"github.com/stretchr/testify/require"
)

func TestConnected(t *testing.T) {
	cases := []struct {
		name                   string
		method                 string
		status                 int
		httpResponse           HTTPResponse
		gatewayConnectedResult bool
	}{
		{
			name:         "405",
			method:       http.MethodPost,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:                   "200 - OK",
			method:                 http.MethodGet,
			status:                 http.StatusOK,
			gatewayConnectedResult: true,
			httpResponse: HTTPResponse{
				Data: true,
			},
		},

		{
			name:                   "200 - OK",
			method:                 http.MethodGet,
			status:                 http.StatusOK,
			gatewayConnectedResult: false,
			httpResponse: HTTPResponse{
				Data: false,
			},
		},
	}

	for _, deviceType := range []deviceWallet.DeviceType{deviceWallet.DeviceTypeUSB, deviceWallet.DeviceTypeEmulator} {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				endpoint := "/connected"
				gateway := &MockGatewayer{}

				if deviceType == deviceWallet.DeviceTypeEmulator {
					endpoint = "/emulator" + endpoint
				}

				gateway.On("Connected").Return(tc.gatewayConnectedResult)

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

					var resp bool
					err = json.Unmarshal(rsp.Data, &resp)
					require.NoError(t, err)

					require.Equal(t, tc.httpResponse.Data.(bool), resp)
				}
			})
		}
	}
}
