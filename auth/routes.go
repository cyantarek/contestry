package auth

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	r.GET("/auth/login", loginHandler)
	r.GET("/logout", logoutHandler)
	r.POST("/auth/login", loginHandler)
	r.GET("/auth/signup", signupHandler)
	r.POST("/auth/signup", signupHandler)
	r.GET("/auth/recovery", recoveryHandler)
	r.POST("/auth/recovery", recoveryHandler)
	r.GET("/auth/confirm/:authToken", confirmationHandler)
}