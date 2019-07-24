package routes

import (
	"github.com/gin-gonic/gin"
	"../controllers"
	"../db"
	"../models"
	"../helper"
)

func GetRouter() *gin.Engine {
	router := gin.New()

	gormDb := db.GetGormDb("contest", "mysql")
	gormDb.Migrate(models.User{}, models.Question{}, models.Contest{}, models.Solution{})

	signKey, verifyKey := helper.SetupAPIKeys()

	h := &controllers.Handler{Store:gormDb, SignKey:signKey, VerifyKey:verifyKey}

	router.Use(gin.Logger())

	router.RedirectTrailingSlash = true
	router.Static("/static", "app/views/static")

	apiRouter := router.Group("/api/v1")

	privilegedRoutes := apiRouter.Group("/privileged", h.JwtValidator())

	authRoutes(apiRouter, h)
	testRoutes(apiRouter, h)
	contestRoutes(privilegedRoutes, h)

	return router
}
