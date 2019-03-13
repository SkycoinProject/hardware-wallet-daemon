package api

import (
	"net/http"

	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
)

// FeaturesResponse is data returned by GET /api/features
type FeaturesResponse struct {
	Features *messages.Features `json:"features"`
}

func features(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.GetFeatures()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
