package router

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/api"
	"github.com/imandaneshi/vite/pkg/config"
)

func InitRoutes(g *gin.Engine){

	// serve frontend static files
	g.Use(static.Serve("/app", static.LocalFile(config.Server.StaticPath, true)))

	g.GET("/:code", api.Redirect)

	g.GET("/", api.HomeRedirect) // redirect `/` to `/app`

	apiRoute := g.Group("api/v1") // api version 1

	// endpoint for shortening the link
	apiRoute.POST("/links", api.CreateShortenLink)
	// user related endpoints
	apiRoute.POST("/users", api.Register)
	apiRoute.POST("/login", api.Login)
}