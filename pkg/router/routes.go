package router

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/api"
	"github.com/imandaneshi/vite/pkg/config"
)

func InitRoutes(g *gin.Engine){
	// serve frontend static files
	g.Use(static.Serve("/", static.LocalFile(config.Server.StaticPath, true)))

	apiRoute := g.Group("api/v1") // api version 1

	// endpoint for shortening the link
	apiRoute.POST("/links", api.CreateShortenLink)
	apiRoute.POST("/users", api.Register)
}