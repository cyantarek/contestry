package user

import (
	"github.com/gin-gonic/gin"
	"../db"
	"fmt"
	"github.com/gin-contrib/sessions"
)

func userProfileHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.HTML(200, "profile", gin.H{"user": user})
}

func judgesListHandler(c *gin.Context) {
	var judges []User
	Db := db.InitDb()
	Db.Find(&judges, "is_judge = true")
	fmt.Println(len(judges))
	c.HTML(200, "judge-list", gin.H{"judges": judges})
}