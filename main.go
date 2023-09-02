package main

import (
	"os"

	rsync "rsync-wrapper/pkg"

	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		logrus.Errorf("error getting current working directory: %v", err)
		os.Exit(1)
	}

	ctx := signals.SetupSignalContext()

	err = rsync.Wrapper(ctx, path)
	if err != nil {
		logrus.Errorf("error from rsync wrapper: %v", err)
		os.Exit(1)
	}
}
