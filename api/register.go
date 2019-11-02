package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
)

type createUser struct {
	Username  string `binding:"required" json:"username" form:"username"`
	Email     string `form:"email" json:"email"`
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Password  string `binding:"required" json:"password" form:"password"`
}

func Register(c *gin.Context) {
	logFields := log.Fields{
		"type":     "api_endpoint",
		"endpoint": "/register",
		"method":   "POST",
	}
	log.WithFields(logFields).Debug("registering a new user")

	var json createUser
	if err := c.ShouldBindJSON(&json); err != nil {
		log.WithFields(logFields).Info("invalid data for registering a new user", err)
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: errors.ValidationError(err.Error(), err)})
		return
	}
	user := &model.User{
		Username:  json.Username,
		Email:     json.Email,
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Password:  json.Password,
	}

	registerError := user.Create()
	if registerError != nil {
		switch registerError.Error() {
		case "already_exists":
			c.AbortWithStatusJSON(400, Response{Ok: false, Error: registerError})
		default:
			c.JSON(500, Response{Ok: false, Error: errors.InternalServerError()})
		}
		return
	}

	user.Password = "" // delete password

	log.WithFields(logFields).Info("successfully registered a new user")
	c.JSON(200, Response{Ok: true, Data: user})
}
