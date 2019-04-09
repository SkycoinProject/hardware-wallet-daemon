package api

import (
	"net/http"
)

// URI: /api/v1/connected
// Method: GET
func connected(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		connected := gateway.Connected()
		writeHTTPResponse(w, HTTPResponse{
			Data: connected,
		})
	}
}
