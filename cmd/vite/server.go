package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/config"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// Server is the cli command that runs our main web server
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Starts the vite web server",
		Before: func(c *cli.Context) error {
			err := model.SetupMongo()
			if err != nil {
				log.Fatal("Error connecting to mongodb", err)
				return err
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mongo-uri",
				Value:       "mongodb://localhost:27017",
				Usage:       "Mongo database uri",
				EnvVars:     []string{"VITE_MONGO_URI"},
				Destination: &config.Database.Uri,
			},
			&cli.StringFlag{
				Name:        "mongo-database",
				Value:       "vite",
				Usage:       "Mongo database name",
				EnvVars:     []string{"VITE_MONGO_DATABASE"},
				Destination: &config.Database.DatabaseName,
			},
			&cli.StringFlag{
				Name:        "static-path",
				Value:       "./client/build",
				Usage:       "Static files path",
				EnvVars:     []string{"VITE_STATIC_PATH"},
				Destination: &config.Server.StaticPath,
			},
			&cli.IntFlag{
				Name:        "server-port",
				Value:       8062,
				Usage:       "Web server port",
				EnvVars:     []string{"VITE_SERVER_PORT", "PORT"},
				Destination: &config.Server.ServerPort,
			},
			&cli.StringFlag{
				Name:        "server-host",
				Value:       "0.0.0.0",
				Usage:       "Web server host",
				EnvVars:     []string{"VITE_SERVER_HOST"},
				Destination: &config.Server.ServerHost,
			},
		},
		Action: func(c *cli.Context) error {
			router := gin.Default()
			// serve frontend static files
			router.Use(static.Serve("/", static.LocalFile(config.Server.StaticPath, true)))
			err := router.Run(fmt.Sprintf("%s:%d", config.Server.ServerHost, config.Server.ServerPort))
			if err != nil {
				log.Fatal("Failed running gin web server", err)
			}
			return nil
		},
	}
}
