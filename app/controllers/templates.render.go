package controllers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Render(ctx *gin.Context, tplName string, ctxData gin.H) {
	//session := sessions.Default(ctx)
	currentUser := h.SessionStore.GetSessionData(ctx, "user")
	if currentUser != nil {
		if ctxData != nil {
			ctx.HTML(200, tplName, gin.H{"current_user": currentUser, "data": ctxData})
		} else {
			ctx.HTML(200, tplName, gin.H{"current_user": currentUser})
		}
	} else {
		ctx.HTML(200, "index", nil)
	}

}
