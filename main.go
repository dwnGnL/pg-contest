package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dwnGnL/pg-contests/internal/cmd"
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/lib/goerrors"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	flagConfig          = "config"
	cliArgMigrationDSN  = "dsn"
	cliArgMigrationDown = "down"
)

var Version = "v0.0.1"

func main() {
	app := &cli.App{
		Name:  "pg-contests",
		Usage: "pg-contests backend",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagConfig,
				Usage:   "set config file",
				Aliases: []string{"c"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "start_service",
				Usage: "start service",
				Action: func(cliContext *cli.Context) error {
					cfg := config.FromFile(cliContext.String(flagConfig))
					fmt.Println(cfg)
					intLogger(cfg.LogLevel)
					return cmd.StartService(cliContext.Context, cfg)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate",
				Action: func(cliContext *cli.Context) error {
					cfg := config.FromFile(cliContext.String(flagConfig))
					fmt.Println(cfg)
					intLogger(cfg.LogLevel)
					return cmd.StartMigrate(cfg)
				},
			},
			{
				Name:  "version",
				Usage: "version",
				Action: func(_ *cli.Context) error {
					fmt.Println(Version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("app.Run: %s", err)
	}
}

func intLogger(logLevel string) {
	var formatter logrus.Formatter = new(logrus.JSONFormatter)
	if os.Getenv("LOG_FORMAT") == "text" {
		formatter = new(logrus.TextFormatter)
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	err = goerrors.Setup(formatter, level)
	if err != nil {
		panic(err)
	}
}
