package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gogo/protobuf/proto"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// TODO(therealssj): add more validation

// TransactionSignRequest is request data for /api/v1/transactionSign
type TransactionSignRequest struct {
	Inputs          []string `json:"inputs"`
	InputIndexes    []uint32 `json:"input_indexes"`
	OutputAddresses []string `json:"output_addresses"`
	Coins           []string `json:"coins"`
	Hours           []string `json:"hours"`
	AddressIndexes  []int    `json:"address_indexes"`
}

// TransactionSignResponse is data returned by POST /api/v1/transactionSign
type TransactionSignResponse struct {
	Signatures *[]string `json:"signatures"`
}

// URI: /api/v1/transactionSign
// Method: POST
// Args: JSON Body
func transactionSign(gateway Gatewayer) http.HandlerFunc {
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

		var req TransactionSignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if err := req.validate(); err != nil {
			logger.WithError(err).Error("invalid sign transaction request")
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		txnInputs, txnOutputs, err := req.TransactionParams()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		msg, err := gateway.TransactionSign(txnInputs, txnOutputs)
		if err != nil {
			logger.Errorf("transactionSign failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

func (r *TransactionSignRequest) validate() error {
	if len(r.Inputs) == 0 {
		return errors.New("inputs is required")
	}

	if len(r.InputIndexes) == 0 {
		return errors.New("input_indexes is required")
	}

	if len(r.Coins) == 0 {
		return errors.New("coins is required")
	}

	if len(r.Hours) == 0 {
		return errors.New("hours is required")
	}

	if len(r.OutputAddresses) == 0 {
		return errors.New("output_addresses is required")
	}

	return nil
}

// TransactionParams returns params for a transaction from the request data
func (r *TransactionSignRequest) TransactionParams() ([]*messages.SkycoinTransactionInput, []*messages.SkycoinTransactionOutput, error) {
	if len(r.Inputs) != len(r.InputIndexes) {
		return nil, nil, errors.New("inputs length not equal to input_indexes length")
	}

	if len(r.OutputAddresses) != len(r.Coins) {
		return nil, nil, errors.New("output_addresses length not equal to coins length")

	}

	if len(r.OutputAddresses) != len(r.Hours) {
		return nil, nil, errors.New("output_addresses length not equal to hours length")
	}

	var transactionInputs []*messages.SkycoinTransactionInput
	var transactionOutputs []*messages.SkycoinTransactionOutput
	for i, input := range r.Inputs {
		var transactionInput messages.SkycoinTransactionInput
		transactionInput.HashIn = proto.String(input)
		transactionInput.Index = proto.Uint32(r.InputIndexes[i])
		transactionInputs = append(transactionInputs, &transactionInput)
	}
	for i, output := range r.OutputAddresses {
		var transactionOutput messages.SkycoinTransactionOutput
		transactionOutput.Address = proto.String(output)

		coins, err := droplet.FromString(r.Coins[i])
		if err != nil {
			return nil, nil, err
		}

		hours, err := strconv.ParseUint(r.Hours[i], 10, 64)
		if err != nil {
			return nil, nil, err
		}

		transactionOutput.Coin = proto.Uint64(coins)
		transactionOutput.Hour = proto.Uint64(hours)
		if i < len(r.AddressIndexes) {
			transactionOutput.AddressIndex = proto.Uint32(uint32(r.AddressIndexes[i]))
		}
		transactionOutputs = append(transactionOutputs, &transactionOutput)
	}

	return transactionInputs, transactionOutputs, nil
}
