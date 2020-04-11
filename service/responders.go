package service

import (
	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, status int, err interface{}) {
	var message interface{}

	switch v := err.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = v
	}

	c.AbortWithStatusJSON(
		status, gin.H{
			"status": status,
			"error":  message,
		},
	)
}

func htmlError(c *gin.Context, err interface{}) {
	c.HTML(400, "error.html", gin.H{"error": err})
	c.Abort()
}

func badRequest(c *gin.Context, err interface{}) {
	errorResponse(c, 400, err)
}
