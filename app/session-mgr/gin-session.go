package session_mgr

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
)

//wrapping concrete dependencies with application specific struct
type GinSession struct {
	session sessions.Store
}

func NewGinSession(name string, secret string, r *gin.Engine) SessionStore {
	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions(name, store))

	return &GinSession{store}
}

func (g *GinSession) GetSessionData(ctx *gin.Context, key string) interface{} {
	session := sessions.Default(ctx)
	data := session.Get(key)
	return data
}

func (g *GinSession) SetSessionData(ctx *gin.Context, key string, data interface{}) {
	session := sessions.Default(ctx)
	session.Set(key, data)
	err := session.Save()
	if err != nil {
		log.Println(err.Error())
	}
}


