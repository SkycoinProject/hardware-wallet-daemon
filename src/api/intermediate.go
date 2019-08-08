package api

import (
	"encoding/json"
	"net/http"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"
)

// PinMatrixRequest request data from /api/v1/intermediate/pin_matrix
type PinMatrixRequest struct {
	Pin string `json:"pin"`
}

func pinMatrixRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req PinMatrixRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		var msg wire.Message
		var err error
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.PinMatrixAck(req.Pin)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				errCH <- 1
				return
			}

			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
		case <-ctx.Done():
			err = gateway.Disconnect()
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// PassPhraseRequest request data from /api/v1/intermediate/passphrase
type PassPhraseRequest struct {
	Passphrase string `json:"passphrase"`
}

func passphraseRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req PassPhraseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		var msg wire.Message
		var err error
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.PassphraseAck(req.Passphrase)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				errCH <- 1
				return
			}

			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
		case <-ctx.Done():
			err = gateway.Disconnect()
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// WordRequest request data from /api/v1/intermediate/word
type WordRequest struct {
	Word string `json:"word"`
}

func wordRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req WordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		defer r.Body.Close()

		var msg wire.Message
		var err error
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.WordAck(req.Word)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				errCH <- 1
				return
			}

			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
		case <-ctx.Done():
			err = gateway.Disconnect()
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func buttonRequestHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var msg wire.Message
		var err error
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.ButtonAck()
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
		case <-ctx.Done():
			err = gateway.Disconnect()
			if err != nil {
				logger.Error(err)
			}
		}
	}
}
