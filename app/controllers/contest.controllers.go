package controllers

import (
	"github.com/gin-contrib/sessions"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"strconv"

	"../models"
	"../utilities"
	"log"
	"fmt"
)


func (h *Handler) ContestList(c *gin.Context) {
	var contests []models.Contest
	h.Store.GetAll(&contests)
	h.Render(c, "contest-list", gin.H{"contests": contests})
}

func (h *Handler) ContestJoinAsTeam(c *gin.Context) {
	contestSlug := c.Param("slug")
	var contest models.Contest

	h.Store.GetSingle(&contest, "slug = ?", contestSlug)

	currentUser := h.SessionStore.GetSessionData(c, "user").(models.User)

	var team models.Team
	h.Store.GetSingle(&team, "leader_id = ?", currentUser.ID)

	h.Store.ExecSQL("insert into team_contest(team_id, contest_id) values(?, ?)", team.ID, contest.ID)

	c.Redirect(http.StatusFound, "/contests/view/"+contest.Slug)
}

func (h *Handler) AddJudgesToContest(c *gin.Context) {
	slug := c.Query("contest")
	var contest models.Contest
	h.Store.GetRaw(&contest, "select * from contests where slug = ?", slug)
	if c.Request.Method == "GET" {
		type ContestJudges struct {
			ID        int
			FirstName string
			LastName  string
		}
		var judges []ContestJudges
		h.Store.GetRaw(&judges,"select users.id, users.first_name, users.last_name from users where users.id not in (select user_id from contest_judges where contest_id = ?) and is_judge = true and is_auth = true", contest.ID)

		h.Render(c, "contest-add-judges", gin.H{"judges": judges, "contest": contest})
	} else if c.Request.Method == "POST" {
		ids := c.PostForm("judges")
		contestID := c.PostForm("contest_id")

		if len(ids) > 1 {
			for _, v := range ids {
				h.Store.ExecSQL("insert into contest_judges (contest_id, user_id) values(?, ?)", contestID, v)
			}
		} else {
			h.Store.ExecSQL("insert into contest_judges (contest_id, user_id) values(?, ?)", contestID, ids)
		}
		c.Redirect(http.StatusFound, "/contests/"+slug)
	}
}

func (h *Handler) ContestDetail(c *gin.Context) {
	slug := c.Param("slug")
	var contest models.Contest
	currentUser := h.SessionStore.GetSessionData(c, "user").(models.User)

	h.Store.GetRaw(&contest, "select * from contests where slug = ?", slug)
	contestTime, err := time.Parse("2006-01-02 00:00:00 +0000 UTC", contest.ContestDate.String())
	if err != nil {
		log.Println(err.Error())
	}
	remaining := contestTime.Sub(time.Now()).Minutes()
	var contestStarted bool
	if remaining < 30.0 {
		contestStarted = true
	}

	if contest.Status == "started" {
		contestStarted = true
	}
	if contest.Status == "active" {
		contestStarted = false
	}

	var count int
	var questionCount int
	var IsAllowed string

	count = h.Store.GetCount("select count(*) from contest_judges where contest_id = ? and user_id = ?", contest.ID, currentUser.ID)
	questionCount = h.Store.GetCount("select count(*) from questions where contest_id = ? and user_id = ?", contest.ID, currentUser.ID)

	type QuestionLimit struct {
		QuestionLimit int
	}
	var qLimit QuestionLimit
	h.Store.GetRaw(&qLimit, "select question_limit from contests where id = ?", contest.ID)

	if count != 0 && questionCount < qLimit.QuestionLimit {
		IsAllowed = "true"
	} else {
		IsAllowed = "false"
	}

	type TeamData struct {
		Name       string
		Slug       string
		LeaderName string
	}
	var teams []TeamData
	h.Store.GetRaw(&teams,"select teams.name, teams.slug,users.username as leader_name from team_contest inner join teams on teams.id = team_contest.team_id inner join users on users.id = teams.leader_id")

	type ContestJudges struct {
		FirstName string
		LastName  string
		Username  string
		Count     int
	}

	var contestJudges []ContestJudges
	h.Store.GetRaw(&contestJudges, "select users.first_name, users.last_name,users.username, count(questions.id) as count from contest_judges left join questions on questions.contest_id = contest_judges.contest_id inner join users on contest_judges.user_id = users.id where contest_judges.contest_id = ? group by users.id", contest.ID)
	fmt.Println(contestJudges)
	type QuestionJudge struct {
		models.Question
		models.User
	}

	var questions []QuestionJudge
	h.Store.GetRaw(&questions, "select * from questions inner join users on questions.user_id = users.id where contest_id = ?", contest.ID)

	var count2 int
	var alreadyJoined bool
	count2 = h.Store.GetCount("select count(*)from team_contest inner join teams on teams.id = team_contest.team_id inner join user_team on user_team.team_id = teams.id where user_team.user_id = ? or teams.leader_id = ? and team_contest.contest_id = ?", currentUser.ID, currentUser.ID, contest.ID)

	if count2 > 0 {
		alreadyJoined = true
	} else {
		alreadyJoined = false
	}

	h.Render(c, "contest-detail", gin.H{"contest": contest, "judges": contestJudges, "questions": questions, "allowed": IsAllowed, "contestStarted": contestStarted, "teams": teams, "alreadyJoined": alreadyJoined})
}

func (h *Handler) CreateNewContest(c *gin.Context) {
	if c.Request.Method == "GET" {
		var allowedJudges []models.User

		h.Store.GetRaw(&allowedJudges, "select * from users where is_judge = true and allow_judge = true")
		h.Render(c, "contest-create", gin.H{"allowedJudges": allowedJudges})
	} else if c.Request.Method == "POST" {
		var contest models.Contest

		contest.Name = c.PostForm("contest-name")
		contest.ContestDate, _ = time.Parse("2006-01-02", c.PostForm("contest-date"))
		contest.RunningTime, _ = strconv.Atoi(c.PostForm("running-time"))
		contest.Slug = utilities.Slugify(contest.Name)
		contest.QuestionLimit, _ = strconv.Atoi(c.PostForm("question-limit"))
		contest.Status = "active"

		h.Store.InsertData(&contest)

		if len(c.PostForm("judges")) > 0 {
			for _, v := range c.PostForm("judges") {
				h.Store.ExecSQL("insert into contest_judges (contest_id, user_id) values(?, ?)", contest.ID, v)
			}
		}

		c.Redirect(http.StatusFound, "/contests")
	}
}

func (h *Handler) AddQuestionToContest(c *gin.Context) {
	if c.Request.Method == "GET" {
		slug := c.Query("contest")
		fmt.Println(slug)
		var contest models.Contest
		h.Store.GetSingle(&contest, "slug = ?", slug)
		h.Render(c, "question-create", gin.H{"contest": contest})
	} else if c.Request.Method == "POST" {
		var question models.Question
		question.Title = c.PostForm("title")
		question.Description = c.PostForm("description")
		question.InputFormat = c.PostForm("input_format")
		question.OutputFormat = c.PostForm("output_format")
		question.Constrains = c.PostForm("constrains")
		question.SampleInput = c.PostForm("sample_input")
		question.SampleOutput = c.PostForm("sample_output")
		question.Input = c.PostForm("input")
		question.CorrectAns = c.PostForm("correct-ans")
		point, _ := strconv.ParseFloat(c.PostForm("point"), 32)
		question.Point = float32(point)
		slug := c.PostForm("slug")
		var contest models.Contest
		h.Store.GetSingle(&contest, "slug = ?", slug)
		fmt.Println(contest)
		question.ContestID = contest.ID
		session := sessions.Default(c)
		userData := session.Get("user").(models.User)
		question.UserID = userData.ID

		h.Store.InsertData(&question)
		c.Redirect(http.StatusFound, "/contests/view/"+slug)
	}
}

func (h *Handler) ViewQuestion(c *gin.Context) {
	questionID := c.Param("id")
	var question models.Question
	var contest models.Contest
	h.Store.GetSingle(&question, "id = ?", questionID)
	h.Store.GetSingle(&contest, "id = ?", question.ContestID)

	currentUser := h.SessionStore.GetSessionData(c, "user").(models.User)
	var solution models.Solution
	h.Store.GetRaw(&solution, "select solutions.id, solutions.result, solutions.point, solutions.exec_time from solutions where solutions.contest_id = ? and solutions.question_id = ? and solutions.user_id = ? order by solutions.id DESC limit 1", contest.ID, questionID, currentUser.ID)
	h.Render(c, "question-detail", gin.H{"question": question, "contest": contest, "output": solution})
}

func (h *Handler) ContestStartByAdmin(c *gin.Context) {
	contestSlug := c.Param("slug")
	var contest models.Contest

	h.Store.GetSingle(&contest, "slug = ?", contestSlug)
	contest.Status = "started"
	h.Store.InsertData(&contest)
	referrer := c.Request.Referer()
	c.Redirect(http.StatusFound, referrer)
}

func (h *Handler) ContestTerminateByAdmin(c *gin.Context) {
	contestSlug := c.Param("slug")
	var contest models.Contest

	h.Store.GetSingle(&contest, "slug = ?", contestSlug)
	contest.Status = "terminated"
	h.Store.InsertData(&contest)
	referrer := c.Request.Referer()
	c.Redirect(http.StatusFound, referrer)
}
