package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/stretchr/testify/require"
)

func TestCancel(t *testing.T) {
	cancelMsg := messages.Failure{
		Code:    messages.FailureType_Failure_ActionCancelled.Enum(),
		Message: newStrPtr("Action canceled by User"),
	}

	msgBytes, err := cancelMsg.Marshal()
	require.NoError(t, err)

	cases := []struct {
		name                string
		method              string
		status              int
		httpResponse        HTTPResponse
		gatewayCancelResult wire.Message
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:   "200 - OK",
			method: http.MethodPut,
			status: http.StatusOK,
			gatewayCancelResult: wire.Message{
				Kind: uint16(messages.MessageType_MessageType_Failure),
				Data: msgBytes,
			},
			httpResponse: HTTPResponse{
				Data: "Action canceled by User",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/cancel"
			gateway := &MockGatewayer{}

			gateway.On("Cancel").Return(tc.gatewayCancelResult, nil)

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
