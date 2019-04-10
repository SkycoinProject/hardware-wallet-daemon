package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/cors"
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

	apiVersion1 = "v1"
)

var (
	logger = logging.MustGetLogger("daemon-api")
)

// corsRegex matches all localhost origin headers
var corsRegex *regexp.Regexp

func init() {
	var err error
	corsRegex, err = regexp.Compile(`^https?://localhost|127\.0\.0\.1:\d+$`)
	if err != nil {
		logger.Panic(err)
	}
}

type muxConfig struct {
	host               string
	enableCSRF         bool
	disableHeaderCheck bool
	hostWhitelist      []string
	mode               deviceWallet.DeviceType
}

// Server exposes an HTTP API
type Server struct {
	server   *http.Server
	listener net.Listener
	done     chan struct{}
}

// Config configures Server
type Config struct {
	EnableCSRF         bool
	DisableHeaderCheck bool
	HostWhitelist      []string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
	Mode               deviceWallet.DeviceType
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

func create(host string, c Config, gateway *Gateway) *Server {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = defaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = defaultIdleTimeout
	}

	mc := muxConfig{
		host:               host,
		enableCSRF:         c.EnableCSRF,
		disableHeaderCheck: c.DisableHeaderCheck,
		hostWhitelist:      c.HostWhitelist,
		mode:               c.Mode,
	}

	srvMux := newServerMux(mc, gateway.Device)

	srv := &http.Server{
		Handler:      srvMux,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		IdleTimeout:  c.IdleTimeout,
	}

	return &Server{
		server: srv,
		done:   make(chan struct{}),
	}
}

// Create create a new http server
func Create(host string, c Config, gateway *Gateway) (*Server, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s := create(host, c, gateway)

	s.listener = listener

	return s, nil
}

func newServerMux(c muxConfig, gateway Gatewayer) *http.ServeMux {
	mux := http.NewServeMux()

	allowedOrigins := []string{
		fmt.Sprintf("http://%s", c.host),
		"https://staging.wallet.skycoin.net",
		"https://wallet.skycoin.net",
	}

	for _, s := range c.hostWhitelist {
		allowedOrigins = append(allowedOrigins, fmt.Sprintf("http://%s", s))
	}

	corsValidator := func(origin string) bool {
		if corsRegex.MatchString(origin) {
			return true
		}

		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == origin {
				return true
			}
		}

		return false
	}

	corsHandler := cors.New(cors.Options{
		AllowOriginFunc:    corsValidator,
		Debug:              false,
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", CSRFHeaderName},
		AllowCredentials:   false, // credentials are not used, but it would be safe to enable if necessary
		OptionsPassthrough: false,
	})

	headerCheck := func(host string, hostWhitelist []string, handler http.Handler) http.Handler {
		handler = originRefererCheck(host, hostWhitelist, handler)
		handler = hostCheck(host, hostWhitelist, handler)
		return handler
	}

	webHandlerWithOptionals := func(endpoint string, handlerFunc http.Handler, checkCSRF, checkHeaders bool) {
		handler := wh.ElapsedHandler(logger, handlerFunc)

		handler = corsHandler.Handler(handler)

		if checkCSRF {
			handler = CSRFCheck(c.enableCSRF, handler)
		}

		if checkHeaders {
			handler = headerCheck(c.host, c.hostWhitelist, handler)
		}

		handler = gziphandler.GzipHandler(handler)
		mux.Handle(endpoint, handler)
	}

	webHandler := func(endpoint string, handler http.Handler) {
		handler = wh.ElapsedHandler(logger, handler)
		webHandlerWithOptionals(endpoint, handler, c.enableCSRF, !c.disableHeaderCheck)
	}

	webHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandler("/api/"+apiVersion1+endpoint, handler)
	}

	// get the current CSRF token
	csrfHandlerV1 := func(endpoint string, handler http.Handler) {
		webHandlerWithOptionals("/api/"+apiVersion1+endpoint, handler, false, !c.disableHeaderCheck)
	}
	csrfHandlerV1("/csrf", getCSRFToken(c.enableCSRF)) // csrf is always available, regardless of the API set

	// hw daemon endpoints
	webHandlerV1("/generate_addresses", generateAddresses(gateway))
	webHandlerV1("/apply_settings", applySettings(gateway))
	webHandlerV1("/backup", backup(gateway))
	webHandlerV1("/cancel", cancel(gateway))
	webHandlerV1("/check_message_signature", checkMessageSignature(gateway))
	webHandlerV1("/features", features(gateway))
	// enable firmware update endpoint only for hw wallet
	if c.mode == deviceWallet.DeviceTypeUSB {
		webHandlerV1("/firmware_update", firmwareUpdate(gateway))
	}
	webHandlerV1("/generate_mnemonic", generateMnemonic(gateway))
	webHandlerV1("/recovery", recovery(gateway))
	webHandlerV1("/set_mnemonic", setMnemonic(gateway))
	webHandlerV1("/set_pin_code", setPinCode(gateway))
	webHandlerV1("/sign_message", signMessage(gateway))
	webHandlerV1("/transaction_sign", transactionSign(gateway))
	webHandlerV1("/wipe", wipe(gateway))
	webHandlerV1("/connected", connected(gateway))

	webHandlerV1("/intermediate/pin_matrix", pinMatrixRequestHandler(gateway))
	webHandlerV1("/intermediate/passphrase", passphraseRequestHandler(gateway))
	webHandlerV1("/intermediate/word", wordRequestHandler(gateway))

	return mux
}

// HandleFirmwareResponseMessages handles response messages from the firmware
func HandleFirmwareResponseMessages(w http.ResponseWriter, r *http.Request, gateway Gatewayer, msg wire.Message) {
	switch msg.Kind {
	case uint16(messages.MessageType_MessageType_PinMatrixRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: "PinMatrixRequest",
		})
	case uint16(messages.MessageType_MessageType_PassphraseRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: "PassPhraseRequest",
		})
	case uint16(messages.MessageType_MessageType_WordRequest):
		writeHTTPResponse(w, HTTPResponse{
			Data: "WordRequest",
		})
	case uint16(messages.MessageType_MessageType_ButtonRequest):
		msg, err := gateway.ButtonAck()
		if err != nil {
			logger.Error(err.Error())
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
			Data: addresses,
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
			Data: features,
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
			Data: signature,
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
			Data: &signatures,
		})
	default:
		resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("recevied unexpected response message type: %s", messages.MessageType(msg.Kind)))
		writeHTTPResponse(w, resp)
	}
}

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

		msg, err := gateway.PinMatrixAck(req.Pin)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
	}
}

// PassPhraseRequest request data from /api/v1/intermediate/passphrases
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

		msg, err := gateway.PassphraseAck(req.Passphrase)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		}

		HandleFirmwareResponseMessages(w, r, gateway, msg)
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

func newUint32Ptr(n uint32) *uint32 {
	return &n
}
