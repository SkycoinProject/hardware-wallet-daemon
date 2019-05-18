package api

import (
	"net/http"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

// URI: /api/v1/wipe
// Method: DELETE
func wipe(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("wipe failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		msg, err := gateway.Wipe()
		if err != nil {
			logger.Errorf("wipe failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, gateway, msg)
	}
}
