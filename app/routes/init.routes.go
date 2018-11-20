package routes

import (
	"github.com/gin-gonic/gin"
	"../views"
	"../controllers"
	"../session-mgr"
	"../db"
	"../models"
	"../api"
)

func GetRouter() *gin.Engine {
	router := gin.New()

	gormDb := db.GetGormDb("contest", "mysql")
	gormDb.Migrate(models.User{}, models.Question{}, models.Team{}, models.Contest{}, models.Solution{})

	//mgoDb := db.GetMongoDb("contest")

	ginSession := session_mgr.NewGinSession("user-session", "secret", router)
	//mockSession := session_mgr.NewMockSession("user-session", "secret", router)

	signKey, verifyKey := api.SetupAPIKeys()

	h := &controllers.Handler{Store:gormDb, SessionStore:ginSession, Type:"template", SignKey:signKey, VerifyKey:verifyKey}

	router.HTMLRender = views.LoadTemplates()
	router.Use(gin.Logger())

	router.RedirectTrailingSlash = true
	router.Static("/static", "app/views/static")

	router.Use(h.TypeChecker())
	authRoutes(router, h)
	contestRoutes(router, h)
	userRoutes(router, h)
	apiRoutes(router, h)
	//profilingRoutes(router)

	return router
}
