package api

import (
	"fmt"
	"net/http"
	"strconv"
)

// URI: /api/v1/recovery
// Method: POST
// Args:
//  word-count: mnemonic seed length
//  use-passphrase: (boolean) ask for passphrase before starting operation
//  dry-run: (bool) perform dry-run recovery workflow (for safe mnemonic validation).
func recovery(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return

		}

		usePassphrase, err := parseBoolFlag(r.FormValue("use-passphrase"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for use-passphrase")
			writeHTTPResponse(w, resp)
			return
		}

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

		dryRun, err := parseBoolFlag(r.FormValue("dry-run"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for dry-run")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.Recovery(uint32(wc), usePassphrase, dryRun)
		if err != nil {
			logger.Errorf("recovery failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}