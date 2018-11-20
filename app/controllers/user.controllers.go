package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"fmt"
)

func (h *Handler) UserProfile(c *gin.Context) {
	h.Render(c, "profile", nil)
}

func (h *Handler) Ranking(c *gin.Context) {
	type Ranking struct {
		Username       string
		Score          int
		TotalContests  int
		TotalQuestions int
	}

	var ranking []Ranking
	h.Store.GetRaw(&ranking, "select distinct users.username, sum(solutions.point) as score, solutions.contest_id as total_contests, count(solutions.question_id) as total_questions from users inner join solutions on solutions.user_id = users.id where users.is_participant = 1 and users.is_auth = 1 group by users.username order by score desc")
	fmt.Println(ranking)
	h.Render(c, "rank-list", gin.H{"ranking": ranking})
}

func (h *Handler) ViewOtherUserProfile(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	h.Store.GetSingle(&user, "username = ?", username)
	h.Render(c, "profile", gin.H{"user": user})
}

func (h *Handler) ListOfJudges(c *gin.Context) {
	var judges []models.User
	h.Store.GetRaw(&judges, "select * from users where is_judge = true")
	h.Render(c, "judge-list", gin.H{"judges": judges})
}
