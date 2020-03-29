package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	skyWallet "github.com/SkycoinProject/hardware-wallet-go/src/skywallet"
	messages "github.com/SkycoinProject/hardware-wallet-protob/go"

	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/gogo/protobuf/proto"
)

// BitcoinTransactionSignRequest is request data for /api/v1/bitcoin_transaction_sign
type BitcoinTransactionSignRequest struct {
	TransactionInputs  []BitcoinTransactionInput  `json:"transaction_inputs"`
	TransactionOutputs []BitcoinTransactionOutput `json:"transaction_outputs"`
}

// BitcoinTransactionInput is a Bitcoin transaction input
type BitcoinTransactionInput struct {
	Index    uint32 `json:"index"`
	PrevHash string `json:"prev_hash"`
}

// BitcoinTransactionOutput is a Bitcoin transaction output
type BitcoinTransactionOutput struct {
	AddressIndex *uint32 `json:"address_index"` // pointer to differentiate between 0 and nil
	Address      string  `json:"address"`
	Coins        string  `json:"coins"`
}

// URI: /api/v1/bitcoin_transaction_sign
// Method: POST
// Args: JSON Body
func bitcoinTransactionSign(gateway Gatewayer) http.HandlerFunc {
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

		var req BitcoinTransactionSignRequest
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

		txnInputs, txnOutputs, err := req.BitcoinTransactionParams()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("transactionSign failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		var signatures []string
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			signer := skyWallet.BitcoinTransactionSigner{
				Inputs:   txnInputs,
				Outputs:  txnOutputs,
				Version:  1,
				LockTime: 0,
			}

			signatures, err = gateway.GeneralTransactionSign(&signer)
			if err != nil {
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			writeHTTPResponse(w, HTTPResponse{
				Data: &signatures,
			})
		case <-errCH:
			logger.Errorf("transactionSign failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		case <-ctx.Done():
			disConnErr := gateway.Disconnect()
			if disConnErr != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				resp := NewHTTPErrorResponse(499, "Client Closed Request")
				writeHTTPResponse(w, resp)
			}
		}
	}
}

func (r *BitcoinTransactionSignRequest) validate() error {
	if len(r.TransactionInputs) == 0 {
		return errors.New("inputs are required")
	}

	for _, input := range r.TransactionInputs {
		if input.PrevHash == "" {
			return errors.New("input hash cannot be empty")
		}
	}

	for _, output := range r.TransactionOutputs {
		if output.Address == "" {
			return errors.New("address cannot be empty")
		}

		if output.Coins == "" {
			return errors.New("coins cannot be empty")
		}

	}

	return nil
}

// BitcoinTransactionParams returns params for a transaction from the request data
func (r *BitcoinTransactionSignRequest) BitcoinTransactionParams() ([]*messages.BitcoinTransactionInput, []*messages.BitcoinTransactionOutput, error) {
	var transactionInputs []*messages.BitcoinTransactionInput
	var transactionOutputs []*messages.BitcoinTransactionOutput

	for _, input := range r.TransactionInputs {
		var transactionInput messages.BitcoinTransactionInput

		decoded, err := hex.DecodeString(input.PrevHash)
		if err != nil {
			return nil, nil, err
		}
		transactionInput.PrevHash = decoded
		transactionInput.AddressN = proto.Uint32(input.Index)
		transactionInputs = append(transactionInputs, &transactionInput)
	}

	for _, output := range r.TransactionOutputs {
		var transactionOutput messages.BitcoinTransactionOutput

		coins, err := droplet.FromString(output.Coins)
		if err != nil {
			return nil, nil, err
		}

		transactionOutput.Address = proto.String(output.Address)
		transactionOutput.Coin = proto.Uint64(coins)

		if output.AddressIndex != nil {
			transactionOutput.AddressIndex = proto.Uint32(*output.AddressIndex)
		}

		transactionOutputs = append(transactionOutputs, &transactionOutput)
	}

	return transactionInputs, transactionOutputs, nil
}
