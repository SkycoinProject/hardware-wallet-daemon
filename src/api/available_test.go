package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAvailable(t *testing.T) {
	cases := []struct {
		name                   string
		method                 string
		status                 int
		gatewayAvailableResult bool
		httpResponse           HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodPost,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:                   "200",
			method:                 http.MethodGet,
			status:                 http.StatusOK,
			gatewayAvailableResult: true,
			httpResponse: HTTPResponse{
				Data: true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/available"
			gateway := &MockGatewayer{}

			gateway.On("Available").Return(tc.gatewayAvailableResult, nil)

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
				var resp bool
				err = json.Unmarshal(rsp.Data, &resp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data, resp)
			}
		})
	}
}
