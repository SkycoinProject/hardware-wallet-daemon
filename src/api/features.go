package api

import (
	"net/http"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

// URI: /api/v1/features
// Method: GET
func features(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("features failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		var msg wire.Message
		var err error
		retCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.GetFeatures()
			if err != nil {
				logger.Errorf("features failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-ctx.Done():
			logger.Error(gateway.Disconnect())
		}
	}
}
