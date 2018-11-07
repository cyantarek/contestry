package contest

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/", indexHandler)
	r.GET("/contests/create", createContestHandler)
	r.GET("/contests", contestListHandler)
	r.POST("/contests/create", createContestHandler)
}