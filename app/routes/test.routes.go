package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func testRoutes(r *gin.RouterGroup, h *controllers.Handler) {
	r.GET("/test", h.TestAuthed)
}