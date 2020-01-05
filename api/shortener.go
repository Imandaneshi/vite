package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// createLink validates request payload for creating a new Link
type createLink struct {
	Address string `binding:"required" json:"address" form:"address"`
}

func CreateShortenLink(c *gin.Context) {
	logFields := log.Fields{
		"type":     "endpoint",
		"endpoint": "/links",
		"method":   "POST",
	}

	var json createLink
	if err := c.ShouldBindJSON(&json); err != nil {
		log.WithFields(logFields).Info("invalid data for creating a shorten link", err)
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: errors.ValidationError(err.Error(), err)})
		return
	}

	user, exists := model.GetUserFromGinContext(c)
	if !exists {
		log.WithFields(logFields).Info("you must be logged in for shorting a link")
		c.AbortWithStatusJSON(http.StatusForbidden, Response{Ok: false, Error: errors.AuthenticationRequired()})
		return
	}

	link, err := model.GenerateRandomShortenLink(json.Address, user.ObjectId)
	if err != nil {
		c.JSON(500, Response{Ok: false, Error: errors.InternalServerError()})
		return
	}

	log.WithFields(logFields).Info("successfully created a shorten link")
	c.JSON(200, Response{Ok: true, Data: link})
}
