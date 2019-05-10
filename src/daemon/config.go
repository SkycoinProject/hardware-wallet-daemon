package daemon

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"

	"github.com/skycoin/skycoin/src/util/file"
)

var (
	help = false
)

// Config records the daemon's configuration
type Config struct {
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string

	// Enable CSRF check
	EnableCSRF bool

	// Disable Host, Origin and Referer header check in the wallet API
	DisableHeaderCheck bool
	// Comma separate list of hostnames to accept in the Host header, used to bypass the Host header check which only applies to localhost addresses
	HostWhitelist string
	hostWhitelist []string

	// Timeouts for the HTTP listener
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	// Logging
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	LogLevel string
	// Enable logging to file
	LogToFile bool

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// Enable HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Expose HTTP profiling on this interface
	HTTPProfHost string

	// Data directory holds app data -- defaults to ~/.skycoin
	DataDirectory string

	// DaemonMode decides with what api is enabled, either wallet or emulator
	DaemonMode string
	daemonMode skyWallet.DeviceType
}

// NewConfig returns a new config instance
func NewConfig(port int, datadir string) Config {
	return Config{
		WebInterfaceAddr: "127.0.0.1",
		WebInterfacePort: port,

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		HTTPReadTimeout:  time.Minute * 10,
		HTTPWriteTimeout: time.Minute * 10,
		HTTPIdleTimeout:  time.Minute * 10,

		// Logging
		ColorLog:  true,
		LogLevel:  "INFO",
		LogToFile: false,

		// disable csrf by default
		EnableCSRF: false,

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",

		// Run daemon in wallet mode by default
		DaemonMode: skyWallet.DeviceTypeUSB.String(),

		DataDirectory: datadir,
	}
}

func (c *Config) postProcess() error {
	if help {
		flag.Usage()
		os.Exit(0)
	}

	var err error
	home := file.UserHome()
	c.DataDirectory, err = file.InitDataDir(replaceHome(c.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.HostWhitelist != "" {
		if c.DisableHeaderCheck {
			return errors.New("host whitelist should be empty when header check is disabled")
		}
		c.hostWhitelist = strings.Split(c.HostWhitelist, ",")
	}

	c.daemonMode = skyWallet.DeviceTypeFromString(c.DaemonMode)
	if c.daemonMode == skyWallet.DeviceTypeInvalid {
		return errors.New("invalid device type")
	}

	return nil
}

// RegisterFlags binds CLI flags to config values
func (c *Config) RegisterFlags() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")
	flag.BoolVar(&c.EnableCSRF, "enable-csrf", c.EnableCSRF, "enable CSRF check")
	flag.BoolVar(&c.DisableHeaderCheck, "disable-header-check", c.DisableHeaderCheck, "disables the host, origin and referer header checks.")
	flag.StringVar(&c.HostWhitelist, "host-whitelist", c.HostWhitelist, "Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.")

	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")

	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	flag.StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")

	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")

	flag.StringVar(&c.DaemonMode, "daemon-mode", c.DaemonMode, "Choices are: USB or EMULATOR")
}

func panicIfError(err error, msg string, args ...interface{}) { // nolint: unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
