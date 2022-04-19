package main

import (
	"fmt"
	"os"

	"github.com/quortex/circleci-docs/pkg/flags"
	log "github.com/sirupsen/logrus"
)

var opts flags.Options

func main() {
	// Flags parsing
	err, exit := flags.Parse(&opts)
	if err != nil {
		log.Error(fmt.Errorf("Invalid flags: %w", err))
		os.Exit(1)
	}
	if exit {
		os.Exit(0)
	}
	log.SetLevel(opts.LogLevel.Level)
}
