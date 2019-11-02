package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
)

type createLink struct {
	Address string `binding:"required" json:"address" form:"address"`
}

func CreateShortenLink(c *gin.Context){
	logFields := log.Fields{
		"type": "api_endpoint",
		"endpoint": "/links",
		"method": "POST",
	}
	log.WithFields(logFields).Debug("creating a new shorten link")

	var json createLink
	if err  := c.ShouldBindJSON(&json); err != nil{
		log.WithFields(logFields).Info("invalid data for creating a shorten link", err)
		c.AbortWithStatusJSON(400, Response{Ok:false, Error: errors.ValidationError(err.Error(), err)})
		return
	}

	link, err := model.GenerateRandomShortenLink(json.Address)
	if err != nil {
		log.WithFields(logFields).Info("internal error creating a shorten link", err)
		c.JSON(500, Response{Ok:false, Error: errors.InternalServerError()})
	}

	log.WithFields(logFields).Info("successfully created a shorten link")
	c.JSON(200, Response{Ok: true, Data:link})
}