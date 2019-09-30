package daemon

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/skycoin/hardware-wallet-daemon/src/api"

	skyWallet "github.com/SkycoinProject/hardware-wallet-go/src/skywallet"

	"github.com/SkycoinProject/skycoin/src/util/file"
)

var (
	help = false
)

// Config records the daemon and build configuration
type Config struct {
	App   AppConfig
	Build api.BuildInfo
}

// AppConfig records the app's configuration
type AppConfig struct {
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

// NewAppConfig returns a new app config instance
func NewAppConfig(port int, datadir string) AppConfig {
	return AppConfig{
		WebInterfaceAddr: "127.0.0.1",
		WebInterfacePort: port,

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
	c.App.DataDirectory, err = file.InitDataDir(replaceHome(c.App.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	if c.App.HostWhitelist != "" {
		if c.App.DisableHeaderCheck {
			return errors.New("host whitelist should be empty when header check is disabled")
		}
		c.App.hostWhitelist = strings.Split(c.App.HostWhitelist, ",")
	}

	c.App.daemonMode = skyWallet.DeviceTypeFromString(c.App.DaemonMode)
	if c.App.daemonMode == skyWallet.DeviceTypeInvalid {
		return errors.New("invalid device type")
	}

	return nil
}

// RegisterFlags binds CLI flags to config values
func (c *AppConfig) RegisterFlags() {
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
