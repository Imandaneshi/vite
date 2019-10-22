package api

import "github.com/gin-gonic/gin"

func LinkShortener(c *gin.Context){
	c.JSON(200, gin.H{"foo": "bar"})
}