package cli

import (
	"os"

	"github.com/pkg/errors"
	"github.com/ulm0/deliver/pkg/deliver"
	"github.com/urfave/cli/v2"
)

// Version show the current version
var Version = "dev"

// Execute executes the actual tool
func Execute() error {
	config := &deliver.Config{}
	app := &cli.App{
		Name:    "deliver",
		Usage:   "A dead simple tool for creating releases (plus artifacts) on Github using GitLab CI",
		Version: Version,
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Pierre Ugaz (ulm0)",
			},
		},
		Flags:  configFlags(config),
		Action: run(config),
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}

func run(c *deliver.Config) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if err := c.Check(); err != nil {
			return errors.Errorf("Failed checking: %s", err)
		}
		if err := c.Execute(); err != nil {
			return errors.Errorf("Failed executing: %s", err)
		}
		return nil
	}
}
