package session_mgr

import "github.com/gin-gonic/gin"

type SessionStore interface {
	GetSessionData(ctx *gin.Context, key string) interface{}
	SetSessionData(ctx *gin.Context, key string, data interface{})
}