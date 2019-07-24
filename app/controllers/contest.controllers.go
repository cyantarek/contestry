package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../models"
	"../utilities"
	"time"
	"fmt"
	"log"
)

func (h *Handler) ApiGetAllContests(c *gin.Context) {
	var contests []models.Contest
	h.Store.GetAll(&contests)
	c.JSON(http.StatusOK, gin.H{"msg": "success", "contests": contests})
}

func (h *Handler) ApiGetAllQuestionsForAContest(c *gin.Context) {
	slug := c.Param("slug")
	var contest models.Contest
	h.Store.GetRaw(&contest, "select * from contests where slug = ?", slug)
	type QuestionJudge struct {
		QuestionId             int    `json:"question_id"`
		QuestionTitle          string `json:"question_title"`
		QuestionSetterUsername string `json:"question_setter_username"`
	}

	var questions []QuestionJudge
	h.Store.GetRaw(&questions, "select questions.id as question_id, questions.title as question_title, users.username as question_setter_username from questions inner join users on questions.user_id = users.id where contest_id = ?", contest.ID)

	c.JSON(http.StatusOK, gin.H{"msg": "success", "contest": contest.Slug, "questions": questions})
}

func (h *Handler) ApiAddQuestionToAContest(c *gin.Context) {
	type payload struct {
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description" binding:"required"`
		InputFormat  string `json:"input_format" binding:"required"`
		OutputFormat string `json:"output_format" binding:"required"`
		Constrains   string `json:"constrains" binding:"required"`
		SampleInput  string `json:"sample_input" binding:"required"`
		SampleOutput string `json:"sample_output" binding:"required"`
		TestInput    string `json:"test_input" binding:"required"`
		TestOutput   string `json:"test_output" binding:"required"`
		TemplateCode string `json:"template_code" binding:"required"`
		Point        float64 `json:"point" binding:"required"`
		ContestID    int `json:"contest_id"`
		UserName     int `json:"user_name"`
	}
	var payloadData payload
	err := c.BindJSON(&payloadData)
	if err != nil {
		log.Println(err.Error())
		return
	}

	count := h.Store.GetCount("select count(*) from questions where title = ?", payloadData.Title)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question already exists"})
		return
	}

	var question models.Question

	question.Title = payloadData.Title
	question.Description = payloadData.Description
	question.InputFormat = payloadData.InputFormat
	question.OutputFormat = payloadData.OutputFormat
	question.Constrains = payloadData.Constrains
	question.SampleInput = payloadData.SampleInput
	question.SampleOutput = payloadData.SampleOutput
	question.TestInput = payloadData.TestInput
	question.TestOutput = payloadData.TestOutput
	question.Point = payloadData.Point
	question.TemplateCode = payloadData.TemplateCode

	slug := c.Param("slug")
	var contest models.Contest
	h.Store.GetSingle(&contest, "slug = ?", slug)

	question.ContestID = contest.ID

	question.UserID = payloadData.UserName

	h.Store.InsertData(&question)

	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
}

func (h *Handler) ApiAddJudgesToAContest(c *gin.Context) {
	type payload struct {
		ContestId int   `json:"contest_id"`
		JudgesId  []int `json:"judges_id"`
	}

	var payloadData payload
	c.BindJSON(&payloadData)

	for _, v := range payloadData.JudgesId {
		h.Store.ExecSQL("insert into contest_judges (contest_id, user_id) values(?, ?)", payloadData.ContestId, v)
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
}

func (h *Handler) ApiCheckAllowedJudgesForAContest(c *gin.Context) {
	slug := c.Param("slug")
	var contest models.Contest
	h.Store.GetRaw(&contest, "select * from contests where slug = ?", slug)
	fmt.Println(contest)
	type ContestJudges struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	var allowedJudges []ContestJudges
	h.Store.GetRaw(&allowedJudges, "select users.id, users.username from users where users.id not in (select user_id from contest_judges where contest_id = ?) and user_type = 'judge' and is_auth = true", contest.ID)

	c.JSON(http.StatusOK, gin.H{"msg": "success", "contest": contest.Slug, "allowed_judges": allowedJudges})
}

func (h *Handler) ApiGetAllJudgesForAContest(c *gin.Context) {
	contestSlug := c.Param("slug")
	var contest models.Contest
	h.Store.GetRaw(&contest, "select * from contests where slug = ?", contestSlug)
	fmt.Println(contest.ID)

	type ContestJudges struct {
		ID            int    `json:"id"`
		Username      string `json:"username"`
		QuestionCount int    `json:"question_count"`
	}

	var contestJudges []ContestJudges
	h.Store.GetRaw(&contestJudges, "select users.id, users.last_name,users.username, count(questions.id) as question_count from contest_judges left join questions on questions.contest_id = contest_judges.contest_id inner join users on contest_judges.user_id = users.id where contest_judges.contest_id = ? group by users.id", contest.ID)

	c.JSON(http.StatusOK, gin.H{"msg": "success", "contest": contestSlug, "judges": contestJudges})
}

func (h *Handler) ApiCreateNewContest(c *gin.Context) {
	fmt.Println(h.ContextData)
	if h.ContextData["user_type"] != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You're not allowed to do this"})
		c.Abort()
		return
	}
	var contest models.Contest

	type payload struct {
		Name          string `json:"name" binding:"required"`
		ContestDate   string `json:"contest_date" binding:"required"`
		RunningTime   int    `json:"running_time" binding:"required"`
		Type          string `json:"type" binding:"required"`
		QuestionLimit int    `json:"question_limit" binding:"required"`
	}

	var payloadData payload

	err := c.BindJSON(&payloadData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "All fields are not filled err: " + err.Error()})
		return
	}

	contest.Name = payloadData.Name
	contest.Type = payloadData.Type
	contest.ContestDate, _ = time.Parse("2006-01-02", payloadData.ContestDate)
	contest.RunningTime = payloadData.RunningTime
	contest.Slug = utilities.Slugify(contest.Name)
	contest.QuestionLimit = payloadData.QuestionLimit
	contest.Status = "active"

	h.Store.InsertData(&contest)
	c.JSON(http.StatusCreated, gin.H{"msg": "new contest created"})
}
