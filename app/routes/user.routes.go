package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func userRoutes(r *gin.Engine, h *controllers.Handler) {
	r.GET("/ranking", h.Authorize(), h.Ranking)
	auth := r.Group("/users", h.Authorize())
	auth.GET("/me", h.UserProfile)
	auth.GET("/judges", h.ListOfJudges)
	auth.GET("/@:username", h.ViewOtherUserProfile)
}
