package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func generateMnemonic(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		// TODO(therealssj): validate word count input for semantic errors?

		wordCount := r.FormValue("word-count")
		if wordCount == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "missing word-count")
			writeHTTPResponse(w, resp)
			return
		}

		wc, err := strconv.ParseUint(wordCount, 10, 32)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid value %s for word-count", wordCount))
			writeHTTPResponse(w, resp)
			return
		}

		usePassphrase, err := parseBoolFlag(r.FormValue("use-passphrase"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for use-passphrase")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.GenerateMnemonic(uint32(wc), usePassphrase)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
