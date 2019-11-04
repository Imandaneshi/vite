package api

import (
	"github.com/gin-gonic/gin"
	"github.com/imandaneshi/vite/pkg/errors"
	"github.com/imandaneshi/vite/pkg/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		token = strings.Replace(token, "Token ", "", 1)
		user, err := model.GetTokenUser(token)
		if err != nil {
			c.Set("user", user)
		}
	}
}

type login struct {
	Username string `binding:"required" json:"username" form:"username"`
	Password string `binding:"required" json:"password" form:"password"`
}

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
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: errors.NotFoundError("User not found")})
		return
	}

	correctPassword := user.ValidatePassword(password)
	if !correctPassword {
		c.AbortWithStatusJSON(400, Response{Ok: false, Error: errors.New("wrong_password", "Password validation failed", nil)})
		return
	} else {
		// create a unique token for user
		token := &model.Token{}
		token.User = user.ObjectId.String()
		err := token.Create()
		if err != nil {
			c.AbortWithStatusJSON(500, Response{Ok: false, Error: errors.InternalServerError()})
			return
		}
		user.Token = token
		c.JSON(200, user)
	}
}
