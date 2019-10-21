package main

import (
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
		Before: func (c *cli.Context) error{
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
		},
		Action: func(c *cli.Context) error {
			r := gin.Default()
			r.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
			r.Run()
			return nil
		},
	}
}
