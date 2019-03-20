package api

import (
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"

	"net/http"
)

// URI: /api/v1/cancel
// Method: PUT
func cancel(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.Cancel()
		if err != nil {
			logger.Errorf("cancel failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
			failureMsg, err := deviceWallet.DecodeFailMsg(msg)
			if err != nil {
				logger.Errorf("cancel failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}

			writeHTTPResponse(w, HTTPResponse{
				Data: failureMsg,
			})
		} else {
			HandleFirmwareResponseMessages(w, r, gateway, msg)
		}
	}
}
