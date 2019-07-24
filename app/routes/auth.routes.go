package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func authRoutes(r *gin.RouterGroup, h *controllers.Handler) {
	authGrp := r.Group("/auth")

	authGrp.POST("/login", h.ApiLoginUser)
	authGrp.POST("/signup", h.ApiSignupUser)
}
