package api

import (
	"net/http"
)

func wipe(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.Wipe()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
