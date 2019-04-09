package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const configuredHost = "127.0.0.1:9510"

func defaultMuxConfig() muxConfig {
	return muxConfig{
		host:        configuredHost,
		disableCSRF: true,
	}
}

var endpointsMethods = map[string][]string{
	"/api/v1/generateAddresses": []string{
		http.MethodPost,
	},
	"/api/v1/applySettings": []string{
		http.MethodPost,
	},
	"/api/v1/backup": []string{
		http.MethodPost,
	},
	"/api/v1/cancel": []string{
		http.MethodPut,
	},
	"/api/v1/checkMessageSignature": []string{
		http.MethodPost,
	},
	"/api/v1/features": []string{
		http.MethodGet,
	},
	"/api/v1/generateMnemonic": []string{
		http.MethodPost,
	},
	"/api/v1/recovery": []string{
		http.MethodPost,
	},
	"/api/v1/setMnemonic": []string{
		http.MethodPost,
	},
	"/api/v1/setPinCode": []string{
		http.MethodPost,
	},
	"/api/v1/signMessage": []string{
		http.MethodPost,
	},
	"/api/v1/transactionSign": []string{
		http.MethodPost,
	},
	"/api/v1/wipe": []string{
		http.MethodDelete,
	},
}

func allEndpoints() []string {
	endpoints := make([]string, len(endpointsMethods))
	i := 0
	for e := range endpointsMethods {
		endpoints[i] = e
		i++
	}
	return endpoints
}

func TestCORS(t *testing.T) {
	cases := []struct {
		name          string
		origin        string
		hostWhitelist []string
		valid         bool
		isHTTPS       bool
	}{
		{
			name:   "options no whitelist",
			origin: configuredHost,
			valid:  true,
		},
		{
			name:   "options no whitelist different localhost port",
			origin: "127.0.0.1:4000",
			valid:  true,
		},

		{
			name:    "options no whitelist skycoin wallet staging site",
			origin:  "staging.wallet.skycoin.net",
			valid:   true,
			isHTTPS: true,
		},

		{
			name:    "options no whitelist skycoin wallet staging site",
			origin:  "wallet.skycoin.net",
			valid:   true,
			isHTTPS: true,
		},

		{
			name:          "options whitelist",
			origin:        "example.com",
			hostWhitelist: []string{"example.com"},
			valid:         true,
		},

		{
			name:   "options no whitelist not whitelisted",
			origin: "example.com",
			valid:  false,
		},

		{
			name:   "options no whitelist check vulnerable domain",
			origin: "127a0a0a1:80",
			valid:  false,
		},
	}

	for _, e := range append(allEndpoints()) {
		for _, tc := range cases {
			for _, m := range []string{http.MethodPost, http.MethodGet} {
				name := fmt.Sprintf("%s %s %s", tc.name, m, e)
				t.Run(name, func(t *testing.T) {
					cfg := defaultMuxConfig()
					cfg.hostWhitelist = tc.hostWhitelist

					req, err := http.NewRequest(http.MethodOptions, e, nil)
					require.NoError(t, err)

					var origin string
					if tc.isHTTPS {
						origin = fmt.Sprintf("https://%s", tc.origin)
					} else {
						origin = fmt.Sprintf("http://%s", tc.origin)
					}

					req.Header.Set("Origin", origin)
					req.Header.Set("Access-Control-Request-Method", m)

					handler := newServerMux(cfg, &MockGatewayer{}, &MockGatewayer{})

					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, req)

					resp := rr.Result()

					allowOrigins := resp.Header.Get("Access-Control-Allow-Origin")
					allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
					allowMethods := resp.Header.Get("Access-Control-Allow-Methods")

					if tc.valid {
						require.Equal(t, origin, allowOrigins)
						require.Equal(t, m, allowMethods)
					} else {
						require.Empty(t, allowOrigins)
						require.Empty(t, allowHeaders)
						require.Empty(t, allowMethods)
					}

					allowCreds := resp.Header.Get("Access-Control-Allow-Credentials")
					require.Empty(t, allowCreds)
				})
			}
		}
	}
}
