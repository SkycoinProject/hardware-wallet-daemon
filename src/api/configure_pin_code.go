package api

import (
	"net/http"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

// URI: /api/v1/configure_pin_code
// Method: POST
// Args:
// - remove_pin: (optional) Used to remove current pin
func configurePinCode(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		removePin, err := parseBoolFlag(r.FormValue("remove_pin"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "invalid value for remove_pin")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("configurePinCode failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.ChangePin(&removePin)
		if err != nil {
			logger.Errorf("configurePinCode failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, gateway, msg)
	}
}
