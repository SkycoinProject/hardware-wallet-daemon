package api

import (
	"net/http"
)

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
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
