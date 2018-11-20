package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func apiRoutes(r *gin.Engine, h *controllers.Handler) {
	apiGrp := r.Group("/api/v1")
	apiGrp.POST("/login", h.UserLogin)
	apiGrp.Use(h.JwtValidator())
	apiGrp.GET("/", h.Index)
}
