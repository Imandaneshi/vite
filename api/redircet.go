package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Redirect(c *gin.Context) {
	code := c.Param("code")

	log.Debugf("looking for link with code %s to redirect user", code)

	link, err := model.GetLink(code)
	if err != nil {
		switch err.Error() {
		case "not_found": // TODO: Redirect to home page
			log.WithFields(log.Fields{"code": code}).Info("link with this code doesn't exists")
			c.AbortWithStatusJSON(404, Response{Ok: false, Error: err})
		default:
			c.JSON(500, Response{Ok: false, Error: errors.InternalServerError()})
		}
		return
	}
	log.WithFields(log.Fields{"code": code}).Infof("redirecting user to destination, %s", link.Address)
	c.Redirect(http.StatusPermanentRedirect, link.Address)
}

func HomeRedirect(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/app")
	log.Info("redirected user to homepage at /app")
}
