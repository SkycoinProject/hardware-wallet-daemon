package api

import (
	"encoding/json"
	"net/http"
)

// GenerateAddressesRequest is request data for /api/v1/generate_addresses
type GenerateAddressesRequest struct {
	AddressN       int  `json:"address_n"`
	StartIndex     int  `json:"start_index"`
	ConfirmAddress bool `json:"confirm_address"`
}

// generateAddresses generates addresses for hardware wallet.
// URI: /api/v1/generate_addresses
// Method: POST
// Args: JSON Body
func generateAddresses(gateway Gatewayer) http.HandlerFunc {
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

		var req GenerateAddressesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		if req.AddressN == 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be 0")
			writeHTTPResponse(w, resp)
			return
		}

		if req.AddressN < 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "address_n cannot be negative")
			writeHTTPResponse(w, resp)
			return
		}

		if req.StartIndex < 0 {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, "start_index cannot be negative")
			writeHTTPResponse(w, resp)
			return
		}

		// simple warning for logs
		if req.AddressN+req.StartIndex > 8 {
			logger.Warnf("wallet generating high index addresses: start_index: %d; address_n: %d", req.StartIndex, req.AddressN)
		}

		msg, err := gateway.AddressGen(req.AddressN, req.StartIndex, req.ConfirmAddress)
		if err != nil {
			logger.Error("generateAddress failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}
