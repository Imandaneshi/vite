package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func DeleteLink(c *gin.Context) {

	logFields := log.Fields{
		"type":     "endpoint",
		"endpoint": "/links",
		"method":   "DELETE",
	}

	code := c.Param("code")
	log.Debugf("looking for link with code %s to delete", code)

	user, exists := model.GetUserFromGinContext(c)
	if !exists {
		log.WithFields(logFields).Info("you must be logged in for deleting a link")
		c.AbortWithStatusJSON(http.StatusForbidden, Response{Ok: false, Error: errors.AuthenticationRequired()})
		return
	}

	link, err := model.GetLink(code)
	if err != nil {
		switch err.Error() {
		case "notFound": // TODO: Redirect to home page
			log.WithFields(log.Fields{"code": code}).Info("Link not found.")
			c.AbortWithStatusJSON(http.StatusNotFound, Response{Ok: false, Error: err})
		default:
			c.JSON(500, Response{Ok: false, Error: errors.InternalServerError()})
		}
		return
	}

	if link.User.Hex() != user.ObjectId.Hex() {
		c.JSON(http.StatusNotFound, Response{Ok: false, Error: errors.NotFoundError("Link not found.")})
		return
	}

	deletingErr := link.Delete()

	if deletingErr != nil {
		log.WithFields(logFields).Infof("Failed deleting link from database: %s", deletingErr.Error())
		c.JSON(http.StatusInternalServerError, Response{Ok: false, Error: errors.InternalServerError()})
		return
	}

	c.Data(204, gin.MIMEJSON, nil)
}