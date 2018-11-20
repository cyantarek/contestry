package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
)

func (h *Handler) Index(c *gin.Context) {
	var contest models.Contest
	//h.Store.GetSingle(&contest, "status = 'active' or status = 'started'", nil)
	h.Store.GetRaw(&contest, "select * from contests where status = 'active' or status = 'started' limit 1")

	if h.Type == "template" {
		h.Render(c, "index", gin.H{"upcomingContest": contest})
	} else if h.Type == "api" {
		c.JSON(200, contest)
	}
}
