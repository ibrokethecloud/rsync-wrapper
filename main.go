package main

import (
	"flag"
	"os"

	rsync "github.com/ibrokethecloud/rsync-wrapper/pkg"
	"github.com/mitchellh/go-homedir"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
)

var (
	verbose bool
)

func main() {
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.Parse()

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	path, err := os.Getwd()
	if err != nil {
		logrus.Errorf("error getting current working directory: %v", err)
		os.Exit(1)
	}

	logrus.Debug("looking up home dir")
	homeDir, err := homedir.Dir()
	if err != nil {
		logrus.Errorf("error getting user home directory: %v", err)
		os.Exit(1)
	}
	ctx := signals.SetupSignalContext()
	logrus.Debugf("home dir %s\n", homeDir)
	err = rsync.Wrapper(ctx, path, homeDir)
	if err != nil {
		logrus.Errorf("error from rsync wrapper: %v", err)
		os.Exit(1)
	}
}
