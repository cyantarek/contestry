package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) TestAuthed(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

