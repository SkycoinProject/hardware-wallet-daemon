package api

import (
	"encoding/json"
	"net/http"
)

// ApplySettingsRequest is request data for /api/v1/apply_settings
type ApplySettingsRequest struct {
	Label         string `json:"label"`
	UsePassphrase bool   `json:"use_passphrase"`
}

// applySettings apply device settings
// URI: /api/v1/apply_settings
// Method: POST
// Args: JSON Body
func applySettings(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req ApplySettingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		msg, err := gateway.ApplySettings(req.UsePassphrase, req.Label)
		if err != nil {
			logger.Error("applySettings failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
