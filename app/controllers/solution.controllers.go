package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../sandbox"
	"../models"
)

/*
func TeamList(c *gin.Context) {
	type Team struct {
		Name        string
		DateCreated time.Time
		Status      string
		Leader      string
		Slug        string
		MemberCount int
	}
	var teams []Team
	store.GetRaw(teams,"select teams.name,teams.date_created,teams.status,users.username as leader,count(user_team.team_id) as member_count,teams.slug from teams inner join users on users.id = teams.leader_id inner join user_team on teams.id = user_team.team_id where teams.status = 'active'")
	session := sessions.Default(c)
	user := session.Get("user").(models.User)
	fmt.Println(user.ID)

	var count int
	var alreadyJoined bool
	count = store.GetCount("select count(*) from user_team where user_team.user_id = ?", user.ID)

	if count > 0 {
		alreadyJoined = true
	} else {
		alreadyJoined = false
	}

	views.Render(c, "team-list", gin.H{"teams": teams, "alreadyJoined": alreadyJoined})
}

func CreateNewTeam(c *gin.Context) {
	if c.Request.Method == "GET" {
		views.Render(c, "team-create", nil)
	} else if c.Request.Method == "POST" {
		var team models.Team

		session := sessions.Default(c)
		currentUser := session.Get("user").(models.User)

		team.LeaderID = currentUser.ID
		team.Name = c.PostForm("team-name")
		team.TagLine = c.PostForm("team-tagline")
		team.Status = "active"
		team.Slug = utilities.Slugify(team.Name)

		store.InsertData(&team)
		c.Redirect(http.StatusFound, "/teams")
	}
}
*/

func (h *Handler) SubmitSolution(c *gin.Context) {
	var newSandbox sandbox.ISandbox
	newSandbox = sandbox.NewSandbox(
		"app/sandbox-store",
		c.PostForm("contest-id"),
		c.PostForm("question-id"),
		c.PostForm("user-id"),
		c.PostForm("code"),
		c.PostForm("language-name"),
	)
	newSandbox.Prepare()
	go newSandbox.Run()

	var solution models.Solution
	h.Store.GetRaw(&solution, "select * from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", c.PostForm("contest-id"), c.PostForm("question-id"), c.PostForm("user-id"))
	if solution.ID > 0 {
		solution.Result = "(update) Solution Submitted, please check after 4-5 seconds"
		solution.Point = -1
		solution.ExecTime = "N/A"
		h.Store.UpdateData(&solution)
		//store.ExecSQL("update solutions set result = '(update) Solution Submitted, please check after 4-5 seconds', exec_time = '0.1ms' where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ?", c.PostForm("contest-id"), c.PostForm("question-id"), c.PostForm("user-id"))
	} else {
		var solution models.Solution
		solution.ContestID = c.PostForm("contest-id")
		solution.QuestionID = c.PostForm("question-id")
		solution.UserID = c.PostForm("user-id")
		solution.Result = "(new) Solution submitted"
		solution.ExecTime = "N/A"
		solution.Point = -1
		h.Store.InsertData(&solution)
	}
	referer := c.Request.Referer()
	c.Redirect(http.StatusFound, referer)
}
