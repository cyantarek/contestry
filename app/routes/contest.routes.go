package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func contestRoutes(r *gin.RouterGroup, h *controllers.Handler) {
	contestGrp := r.Group("/contests")

	contestGrp.GET("/", h.ApiGetAllContests)
	contestGrp.POST("/", h.AdminValidator(), h.ApiCreateNewContest)
	contestGrp.GET("/:slug/allowed-judges", h.AdminValidator(), h.ApiCheckAllowedJudgesForAContest)
	contestGrp.GET("/:slug/judges", h.AdminValidator(), h.ApiGetAllJudgesForAContest)
	contestGrp.POST("/:slug/judges", h.AdminValidator(), h.ApiAddJudgesToAContest)
	contestGrp.GET("/:slug/questions", h.AdminValidator(), h.ApiGetAllQuestionsForAContest)
	contestGrp.POST("/:slug/questions", h.JudgeValidator(), h.ApiAddQuestionToAContest)
}
