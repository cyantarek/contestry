package contest

import (
	"github.com/gin-gonic/gin"
	"../db"
	"log"
	"time"
	"fmt"
	"github.com/gin-contrib/sessions"
	"net/http"
)

func indexHandler(c *gin.Context) {
	session := sessions.Default(c)
	currentUser := session.Get("user")
	if currentUser == nil {
		http.Redirect(c.Writer, c.Request, "/auth/login", http.StatusTemporaryRedirect)
		return
	}
	c.HTML(200, "index", gin.H{"current_user": currentUser})
}

func contestListHandler(c *gin.Context) {
	var contests []Contest
	Db := db.InitDb()
	Db.Find(&contests)
	c.HTML(200, "contest-list", gin.H{"contests": contests})
}

func createContestHandler(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(200, "contest-create", nil)
	} else if c.Request.Method == "POST" {
		var payload Contest
		err := c.ShouldBind(&payload)
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(payload.ContestDate)
		parsedTime, err := time.Parse("2006-01-02", c.PostForm("contest-date"))
		if err != nil {
			log.Println(err.Error())
			return
		}
		payload.ContestDate = parsedTime

		Db := db.InitDb()
		Db.Create(&payload)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

}