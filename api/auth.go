package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

// AuthMiddleware Checks if token exists in Authorization header or authToken cookie
// and adds User object to gin Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// check if token exists in headers

		log.Debug("Looking for Authorization Header")
		token := c.GetHeader("Authorization")
		if token != "" {
			token = strings.Replace(token, "Token ", "", 1)
			user, err := model.GetTokenUser(token)
			if err == nil {
				log.WithFields(log.Fields{"user": user}).Info("Found user using provided token")
				c.Set("user", user)
				return
			}

		} else {

			// if Authorization header doesn't exists, look for authToken cookie

			log.Debug("Looking for authToken cookie")
			cookie, err := c.Cookie("authToken")
			if err == nil {
				user, err := model.GetTokenUser(cookie)
				if err != nil {
					log.WithFields(log.Fields{"user": user}).Info("Found user using provided token in authToken")
					c.Set("user", user)
					return
				}

			}
		}

	}
}

// login Validates auth endpoint request payload
type login struct {
	Username string `binding:"required" json:"username" form:"username"`
	Password string `binding:"required" json:"password" form:"password"`
}

// Login logs user in and creates a unique token for user
func Login(c *gin.Context) {

	var json login
	if err := c.ShouldBindJSON(&json); err != nil {
		log.Info("invalid data for login", err)
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: errors.ValidationError(err.Error(), err)})
		return
	}

	username := json.Username
	password := json.Password
	user, err := model.GetUserByUsername(username)
	if err != nil {
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: userNotFoundError})
		return
	}

	correctPassword := user.ValidatePassword(password)
	if !correctPassword {
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: wrongPasswordError})
		return
	} else {
		// create a unique token for user
		token := &model.Token{}
		token.User = user.ObjectId
		err := token.Create()
		if err != nil {
			c.AbortWithStatusJSON(500, Response{Ok: false, Error: errors.InternalServerError()})
			return
		}
		user.Token = token
		c.JSON(200, Response{true, user, nil})
	}

}

var (
	wrongPasswordError = errors.New("wrongPassword", "Password validation failed", nil)
	userNotFoundError  = errors.NotFoundError("User not found")
)
