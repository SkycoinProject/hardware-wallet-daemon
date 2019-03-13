package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	deviceWallet "github.com/skycoin/hardware-wallet-go/src/device-wallet"
	messages "github.com/skycoin/hardware-wallet-go/src/device-wallet/messages/go"
	"github.com/skycoin/hardware-wallet-go/src/device-wallet/wire"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	defaultReadTimeout  = time.Second * 10
	defaultWriteTimeout = time.Second * 60
	defaultIdleTimeout  = time.Second * 120

	// ContentTypeJSON json content type header
	ContentTypeJSON = "application/json"
	// ContentTypeForm form data content type header
	ContentTypeForm = "application/x-www-form-urlencoded"
)

var (
	logger = logging.MustGetLogger("daemon-api")
)

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// HTTPResponse represents the http response struct
type HTTPResponse struct {
	Error *HTTPError  `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// ReceivedHTTPResponse parsed is a Parsed HTTPResponse
type ReceivedHTTPResponse struct {
	Error *HTTPError      `json:"error,omitempty"`
	Data  json.RawMessage `json:"data"`
}

// HTTPError is included in an HTTPResponse
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NewHTTPErrorResponse returns an HTTPResponse with the Error field populated
func NewHTTPErrorResponse(code int, msg string) HTTPResponse {
	if msg == "" {
		msg = http.StatusText(code)
	}

	return HTTPResponse{
		Error: &HTTPError{
			Code:    code,
			Message: msg,
		},
	}
}

func writeHTTPResponse(w http.ResponseWriter, resp HTTPResponse) {
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		wh.Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", ContentTypeJSON)

	if resp.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if resp.Error.Code < 400 || resp.Error.Code >= 600 {
			logger.Critical().Errorf("writeHTTPResponse invalid error status code: %d", resp.Error.Code)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(resp.Error.Code)
		}
	}

	if _, err := w.Write(out); err != nil {
		logger.WithError(err).Error("http Write failed")
	}
}

// Serve serves the web interface on the configured host
func (s *Server) Serve() error {
	defer close(s.done)

	if err := s.server.Serve(s.listener); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}
	return nil
}

// Shutdown closes the HTTP service. This can only be called after Serve or ServeHTTPS has been called.
func (s *Server) Shutdown() {
	if s == nil {
		return
	}

	logger.Info("Shutting down web interface")
	defer logger.Info("Web interface shut down")
	if err := s.listener.Close(); err != nil {
		logger.WithError(err).Warning("s.listener.Close() error")
	}
	<-s.done
}

func create(host string, c Config, gateway *Gateway) (*Server, error) {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaultIdleTimeout
	}

	srvMux := newServerMux(gateway.USBDevice, gateway.EmulatorDevice)

	srv := &http.Server{
		Handler:      srvMux,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	return &Server{
		server: srv,
		done:   make(chan struct{}),
	}, nil
}

func Create(host string, c Config, gateway *Gateway) (*Server, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s, err := create(host, c, gateway)
	if err != nil {
		if closeErr := s.listener.Close(); closeErr != nil {
			logger.WithError(err).Warning("s.listener.Close() error")
		}
		return nil, err
	}

	s.listener = listener

	return s, nil
}

func newServerMux(usbGateway, emulatorGateway Gatewayer) *http.ServeMux {
	mux := http.NewServeMux()

	webHandler := func(endpoint string, handler http.Handler) {
		handler = wh.ElapsedHandler(logger, handler)

		mux.Handle("/api"+endpoint, handler)
	}

	// hw wallet endpoints
	webHandler("/generateAddresses", generateAddresses(usbGateway))
	webHandler("/applySettings", applySettings(usbGateway))
	webHandler("/backup", backup(usbGateway))
	webHandler("/cancel", cancel(usbGateway))
	webHandler("/checkMessageSignature", checkMessageSignature(usbGateway))
	webHandler("/features", features(usbGateway))
	webHandler("/generateMnemonic", generateMnemonic(usbGateway))
	webHandler("/recovery", recovery(usbGateway))
	webHandler("/setMnemonic", setMnemonic(usbGateway))
	webHandler("/setPinCode", setPinCode(usbGateway))
	webHandler("/signMessage", signMessage(usbGateway))
	webHandler("/transactionSign", transactionSign(usbGateway))
	webHandler("/wipe", wipe(usbGateway))
	webHandler("/intermediate/pinmatrix", PinMatrixRequestHandler(usbGateway))
	webHandler("/intermediate/passphrase", PassphraseRequestHandler(usbGateway))
	webHandler("/intermediate/word", WordRequestHandler(usbGateway))

	// emulator endpoints
	webHandler("/emulator/generateAddresses", generateAddresses(emulatorGateway))
	webHandler("/emulator/applySettings", applySettings(emulatorGateway))
	webHandler("/emulator/backup", backup(emulatorGateway))
	webHandler("/emulator/cancel", cancel(emulatorGateway))
	webHandler("/emulator/checkMessageSignature", checkMessageSignature(emulatorGateway))
	webHandler("/emulator/features", features(emulatorGateway))
	webHandler("/emulator/generateMnemonic", generateMnemonic(emulatorGateway))
	webHandler("/emulator/recovery", recovery(emulatorGateway))
	webHandler("/emulator/setMnemonic", setMnemonic(emulatorGateway))
	webHandler("/emulator/setPinCode", setPinCode(emulatorGateway))
	webHandler("/emulator/signMessage", signMessage(emulatorGateway))
	webHandler("/emulator/transactionSign", transactionSign(emulatorGateway))
	webHandler("/emulator/wipe", wipe(emulatorGateway))
	webHandler("/emulator/intermediate/pinmatrix", PinMatrixRequestHandler(emulatorGateway))
	webHandler("/emulator/intermediate/passphrase", PassphraseRequestHandler(emulatorGateway))
	webHandler("/emulator/intermediate/word", WordRequestHandler(emulatorGateway))

	return mux
}

func parseBoolFlag(v string) (bool, error) {
	if v == "" {
		return false, nil
	}

	return strconv.ParseBool(v)
}

type IntermediateResponse struct {
	RequestType string `json:"request_type"`
}

func HandleFirmwareResponseMessages(w http.ResponseWriter, r *http.Request, gateway Gatewayer, msg wire.Message) {
	switch msg.Kind {
	case uint16(messages.MessageType_MessageType_PinMatrixRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "PinMatrixRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_PassphraseRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "PassPhraseRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_WordRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: IntermediateResponse{
				RequestType: "WordRequest",
			},
		})
	case uint16(messages.MessageType_MessageType_ButtonRequest):
		msg, err := gateway.ButtonAck()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnauthorized, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	case uint16(messages.MessageType_MessageType_Failure):
		failureMsg, err := deviceWallet.DecodeFailMsg(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}
		resp := NewHTTPErrorResponse(http.StatusConflict, failureMsg)
		writeHTTPResponse(w, resp)
		return
	case uint16(messages.MessageType_MessageType_Success):
		successMsg, err := deviceWallet.DecodeSuccessMsg(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnauthorized, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: successMsg,
		})
	// AddressGen Response
	case uint16(messages.MessageType_MessageType_ResponseSkycoinAddress):
		addresses, err := deviceWallet.DecodeResponseSkycoinAddress(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: GenerateAddressesResponse{
				Addresses: addresses,
			},
		})
	// Features Response
	case uint16(messages.MessageType_MessageType_Features):
		features := &messages.Features{}
		err := proto.Unmarshal(msg.Data, features)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: FeaturesResponse{
				Features: features,
			},
		})
	// SignMessage Response
	case uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage):
		signature, err := deviceWallet.DecodeResponseSkycoinSignMessage(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: SignMessageResponse{
				Signature: signature,
			},
		})
	// TransactionSign Response
	case uint16(messages.MessageType_MessageType_ResponseTransactionSign):
		signatures, err := deviceWallet.DecodeResponseTransactionSign(msg)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: TransactionSignResponse{
				Signatures: signatures,
			},
		})
	default:
		resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("recevied unexpected response message type: %s", messages.MessageType(msg.Kind)))
		writeHTTPResponse(w, resp)
	}
}

type PinMatrixRequest struct {
	Pin string `json:"pin"`
}

func PinMatrixRequestHandler(gateway Gatewayer) http.HandlerFunc {
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

		msg, err := gateway.PinMatrixAck(req.Pin)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

type PassPhraseRequest struct {
	Passphrase string `json:"passphrase"`
}

func PassphraseRequestHandler(gateway Gatewayer) http.HandlerFunc {
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

		msg, err := gateway.PassphraseAck(req.Passphrase)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

type WordRequest struct {
	Word string `json:"word"`
}

func WordRequestHandler(gateway Gatewayer) http.HandlerFunc {
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

		msg, err := gateway.WordAck(req.Word)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

func newStrPtr(s string) *string {
	return &s
}
