package api

import (
	"net/http"

	"github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
)

// FeaturesResponse is data returned by GET /api/v1/features
type FeaturesResponse struct {
	Features *messages.Features `json:"features"`
}

// URI: /api/v1/features
// Method: GET
func features(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.GetFeatures()
		if err != nil {
			logger.Errorf("features failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
