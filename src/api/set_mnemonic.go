package api

import (
	"net/http"
)

// URI: /api/v1/setMnemonic
// Method: POST
// Args:
//  mnemonic: bip39 mnemonic seed [required]
func setMnemonic(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		mnemonic := r.FormValue("mnemonic")
		if mnemonic == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "mnemonic is required")
			writeHTTPResponse(w, resp)
			return
		}

		// TODO(therealssj): add mnemonic check?

		msg, err := gateway.SetMnemonic(mnemonic)
		if err != nil {
			logger.Errorf("setMnemonic failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
