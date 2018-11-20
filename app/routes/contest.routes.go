package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
)

func contestRoutes(r *gin.Engine, h *controllers.Handler) {
	r.GET("/", h.Authorize(), h.Index)

	contestGrp := r.Group("/contests", h.Authorize())
	contestGrp.GET("/", h.ContestList)
	contestGrp.GET("/view/:slug", h.ContestDetail)
	contestGrp.GET("/view/:slug/questions/:id", h.ViewQuestion)
	contestGrp.GET("/create", h.CreateNewContest)
	contestGrp.GET("/add-judges", h.AddJudgesToContest)
	contestGrp.GET("/view/:slug/start", h.ContestStartByAdmin)
	contestGrp.GET("/view/:slug/terminate", h.ContestTerminateByAdmin)
	contestGrp.POST("/add-judges", h.AddJudgesToContest)
	contestGrp.POST("/create", h.CreateNewContest)
	//contestGrp.GET("/teams/join/:slug", controllers.ContestJoinAsTeam)

	questionGrp := r.Group("/questions", h.Authorize())
	questionGrp.GET("/create", h.AddQuestionToContest)
	questionGrp.GET("/view/:id", h.ViewQuestion)
	questionGrp.POST("/create", h.AddQuestionToContest)

	solutionGrp := r.Group("/solutions", h.Authorize())
	solutionGrp.POST("/submit", h.SubmitSolution)

	//teamGrp := r.Group("/teams", middlewares.Authorize())
	//teamGrp.GET("/", controllers.TeamList)
	//teamGrp.GET("/create", controllers.CreateNewTeam)
	//teamGrp.POST("/create", controllers.CreateNewTeam)
}
