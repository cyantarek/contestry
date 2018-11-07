package user

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	r.GET("/users/profile", userProfileHandler)
	r.GET("/users/judges",judgesListHandler)
}