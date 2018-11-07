package router

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions/cookie"

	"../templates"
	"../user"
	"../auth"
	"../contest"
	"github.com/gin-contrib/sessions"
)

func GetRouter() http.Handler {
	router := gin.New()
	router.HTMLRender = templates.LoadTemplates()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("user-session", store))
	router.RedirectTrailingSlash = true

	user.Routes(router)
	auth.Routes(router)
	contest.Routes(router)

	return router
}


