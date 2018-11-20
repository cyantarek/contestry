package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func authRoutes(r *gin.Engine, h *controllers.Handler) {
	r.GET("/auth/login", h.UnAuthorize(), h.UserLogin)
	r.GET("/auth/logout", h.Authorize(), h.UserLogout)
	r.POST("/auth/login", h.UnAuthorize(), h.UserLogin)
	r.GET("/auth/signup", h.UnAuthorize(), h.UserSignup)
	r.POST("/auth/signup", h.UnAuthorize(), h.UserSignup)
	//r.GET("/auth/recovery", h.AccountRecovery)
	//r.POST("/auth/recovery", h.AccountRecovery)
	r.GET("/auth/confirm/:authToken", h.UnAuthorize(), h.ConfirmationManager)
}
