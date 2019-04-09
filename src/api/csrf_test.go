package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sync"

	"github.com/stretchr/testify/require"
)

const (
	tokenValid            = "token_valid"
	tokenInvalid          = "token_invalid"
	tokenInvalidSignature = "token_invalid_signature"
	tokenExpired          = "token_expired"
	tokenEmpty            = "token_empty"
)

func setCSRFParameters(t *testing.T, tokenType string, req *http.Request) {
	token, err := newCSRFToken()
	require.NoError(t, err)
	// token check
	switch tokenType {
	case tokenValid:
		req.Header.Set("X-CSRF-Token", token)
	case tokenInvalid:
		// add invalid token value
		req.Header.Set("X-CSRF-Token", "xcasadsadsa")
	case tokenInvalidSignature:
		req.Header.Set("X-CSRF-Token", "YXNkc2Fkcw.YXNkc2Fkcw")
	case tokenExpired:
		// set some old unix time
		expiredToken, err := newCSRFTokenWithTime(time.Unix(1517509381, 10))
		require.NoError(t, err)
		req.Header.Set("X-CSRF-Token", expiredToken)
	case tokenEmpty:
		// add empty token
		req.Header.Set("X-CSRF-Token", "")
	}
}

func TestCSRFWrapper(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	cases := []string{tokenInvalid, tokenExpired, tokenEmpty, tokenInvalidSignature}

	for endpoint := range endpointsMethods {
		for _, method := range methods {
			for _, c := range cases {
				name := fmt.Sprintf("%s %s %s", method, endpoint, c)
				t.Run(name, func(t *testing.T) {
					gateway := &MockGatewayer{}

					req, err := http.NewRequest(method, endpoint, nil)
					require.NoError(t, err)

					setCSRFParameters(t, c, req)

					rr := httptest.NewRecorder()

					cfg := defaultMuxConfig()
					cfg.enableCSRF = true
					handler := newServerMux(cfg, gateway, gateway)

					handler.ServeHTTP(rr, req)

					status := rr.Code
					require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)

					var errMsg error
					switch c {
					case tokenInvalid, tokenEmpty:
						errMsg = ErrCSRFInvalid
					case tokenInvalidSignature:
						errMsg = ErrCSRFInvalidSignature
					case tokenExpired:
						errMsg = ErrCSRFExpired
					}

					require.Equal(t, fmt.Sprintf("{\n    \"error\": {\n        \"message\": \"%s\",\n        \"code\": 403\n    }\n}", errMsg), rr.Body.String())
				})
			}
		}
	}
}

func TestCSRFWrapperConcurrent(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	cases := []string{tokenInvalid, tokenExpired, tokenEmpty, tokenInvalidSignature}

	gateway := &MockGatewayer{}

	cfg := defaultMuxConfig()
	cfg.enableCSRF = true
	handler := newServerMux(cfg, gateway, gateway)

	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for endpoint := range endpointsMethods {
				for _, method := range methods {
					for _, c := range cases {
						name := fmt.Sprintf("%s %s %s", method, endpoint, c)
						t.Run(name, func(t *testing.T) {

							req, err := http.NewRequest(method, endpoint, nil)
							require.NoError(t, err)

							setCSRFParameters(t, c, req)

							rr := httptest.NewRecorder()

							handler.ServeHTTP(rr, req)

							status := rr.Code
							require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)

							var errMsg error
							switch c {
							case tokenInvalid, tokenEmpty:
								errMsg = ErrCSRFInvalid
							case tokenInvalidSignature:
								errMsg = ErrCSRFInvalidSignature
							case tokenExpired:
								errMsg = ErrCSRFExpired
							}

							require.Equal(t, fmt.Sprintf("{\n    \"error\": {\n        \"message\": \"%s\",\n        \"code\": 403\n    }\n}", errMsg), rr.Body.String())
						})
					}
				}
			}
		}()
	}
	wg.Wait()

}
