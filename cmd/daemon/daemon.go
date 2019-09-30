package main

import (
	"flag"
	"os"

	"github.com/skycoin/hardware-wallet-daemon/src/api"

	"github.com/SkycoinProject/skycoin/src/util/logging"

	"github.com/skycoin/hardware-wallet-daemon/src/daemon"
)

var (
	// Version of the node. Can be set by -ldflags
	Version = "0.1.0"
	// Commit ID. Can be set by -ldflags
	Commit = ""
	// Branch name. Can be set by -ldflags
	Branch = ""

	logger = logging.MustGetLogger("hw-daemon")

	appConfig = daemon.NewAppConfig(
		9510,
		"$HOME/.skycoin")

	parseFlags = true
)

func init() {
	appConfig.RegisterFlags()
}

func main() {
	if parseFlags {
		flag.Parse()
	}

	d := daemon.NewDaemon(daemon.Config{
		App: appConfig,
		Build: api.BuildInfo{
			Version: Version,
			Commit:  Commit,
			Branch:  Branch,
		},
	}, logger)

	// parse config values
	if err := d.ParseConfig(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	if err := d.Run(); err != nil {
		os.Exit(1)
	}
}
