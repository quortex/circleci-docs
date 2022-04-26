package main

import (
	"fmt"
	"os"

	"github.com/quortex/circleci-docs/pkg/circleci"
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

	// Configuration parsing
	log.Debugf("Parsing circleci configuration: %s", opts.Positional.ConfigFile)
	c, err := circleci.NewConfig(opts.Positional.ConfigFile)
	if err != nil {
		log.Error(fmt.Errorf("Cannot parse project: %w", err))
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("Parsed config: %+v", c))
}
