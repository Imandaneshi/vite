package main

import (
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
		},
		Action: func(c *cli.Context) error {
			router := gin.Default()

			// serve frontend static files
			router.Use(static.Serve("/", static.LocalFile(config.Server.StaticPath, true)))

			return nil
		},
	}
}
