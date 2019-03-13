package daemon

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
)

var (
	help = false
)

type Config struct {
	// Remote web interface port
	WebInterfacePort int
	// Remote web interface address
	WebInterfaceAddr string

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
}

func NewConfig(port int, datadir string) Config {
	return Config{
		WebInterfacePort: port,
		WebInterfaceAddr: "127.0.0.1",

		// Timeout settings for http.Server
		// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		HTTPReadTimeout:  time.Second * 10,
		HTTPWriteTimeout: time.Second * 60,
		HTTPIdleTimeout:  time.Second * 120,

		// Logging
		ColorLog:  true,
		LogLevel:  "INFO",
		LogToFile: false,

		// Enable cpu profiling
		ProfileCPU: false,
		// Where the file is written to
		ProfileCPUFile: "cpu.prof",
		// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
		HTTPProf:     false,
		HTTPProfHost: "localhost:6060",

		DataDirectory: datadir,
	}
}

func (c *Config) postProcess() error {
	if help {
		flag.Usage()
	}

	var err error
	home := file.UserHome()
	c.DataDirectory, err = file.InitDataDir(replaceHome(c.DataDirectory, home))
	panicIfError(err, "Invalid DataDirectory")

	return nil
}

func (c *Config) RegisterFlags() {
	flag.BoolVar(&help, "help", false, "Show help")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr", c.WebInterfaceAddr, "addr to serve web interface on")

	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog, "Add terminal colors to log output")
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Choices are: debug, info, warn, error, fatal, panic")
	flag.BoolVar(&c.LogToFile, "logtofile", c.LogToFile, "log to file")

	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU, "enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file", c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf, "run the HTTP profiling interface")
	flag.StringVar(&c.HTTPProfHost, "http-prof-host", c.HTTPProfHost, "hostname to bind the HTTP profiling interface to")

	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory, "directory to store app data (defaults to ~/.skycoin)")
}

func panicIfError(err error, msg string, args ...interface{}) { // nolint: unparam
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func replaceHome(path, home string) string {
	return strings.Replace(path, "$HOME", home, 1)
}
