package api

import (
	"net/http"
)

// applySettings apply device settings
// URI: /api/v1/applySettings
// Method: POST
// Args:
//  label: label for hardware wallet
//  use-passphrase: (boolean) ask for passphrase before starting operation
func applySettings(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		passphrase, err := parseBoolFlag(r.FormValue("use-passphrase"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for use-passphrase")
			writeHTTPResponse(w, resp)
			return
		}

		label := r.FormValue("label")
		if label == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "missing label")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.ApplySettings(passphrase, label)
		if err != nil {
			logger.Error("applySettings failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
