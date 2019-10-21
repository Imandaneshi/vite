package main

import (
	"github.com/imandaneshi/vite/pkg/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
	"os"
	"strings"
	"time"
)


// setupLogging sets logging level for logrus
func setupLogging() {
	switch strings.ToLower(config.Logging.Level) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}


// cliFlags returns global cli flags
func cliFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Value:       true,
			Usage:       "Activate debug information",
			EnvVars:     []string{"VITE_DEBUG"},
			Destination: &config.Server.Debug,
		},
		&cli.StringFlag{
			Name:        "logging-level",
			Value:       "info",
			Usage:       "set logging level",
			EnvVars:     []string{"VITE_LOG_LEVEL"},
			Destination: &config.Logging.Level,
		},
	}
}

func main() {
	app := &cli.App{
		Name:      "Vite",
		Usage:     "vite web server",
		Compiled:  time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Iman Daneshi",
				Email: "emandaneshikohan@gmail.com",
			},
		},
		Flags:    cliFlags(),
		Commands: []*cli.Command{
			Server(),
		},
		Before: func(c *cli.Context) error {
			setupLogging()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("failed starting the web server")
	}
}
