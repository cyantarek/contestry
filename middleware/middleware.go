package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"fmt"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		loginUser := session.Get("user")
		if loginUser == nil {
			fmt.Println("Called")
		}
	}
}
