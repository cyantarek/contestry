package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/smtp"
	"log"
	"regexp"
	"../user"
	"../db"
	"golang.org/x/crypto/bcrypt"
	"time"
	"fmt"
	"github.com/gin-contrib/sessions"
)

func loginHandler(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var payload LoginPayload
		var userCheck user.User

		c.ShouldBind(&payload)

		Db := db.InitDb()
		Db.First(&userCheck, "username = ?", payload.Username)
		if userCheck.Username == "" {
			c.HTML(200, "login", gin.H{})
			return
		} else if !userCheck.IsAuth {
			c.HTML(200, "login", gin.H{})
			return
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(userCheck.Password), []byte(payload.Password))
			if err != nil {
				c.HTML(200, "login", gin.H{})
				return
			} else {
				userCheck.LastLogin = time.Now()
				userCheck.IsLoggedIn = true
				Db.Save(&userCheck)
				session := sessions.Default(c)
				session.Set("user", userCheck)
				err := session.Save()
				if err != nil {
					log.Println(err.Error())
				}
				c.Redirect(301, "/")
				return
			}
		}
	}
	session := sessions.Default(c)
	userL := session.Get("user")
	fmt.Println(userL)
	c.HTML(200, "login", gin.H{"success": "true"})
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("user", nil)
	session.Save()
	fmt.Println("Called")
	c.Redirect(308, "/auth/login")
}

func signupHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(200, "signup", gin.H{"messages": ""})
	} else if c.Request.Method == http.MethodPost {
		var payload RegistrationPayload
		var errMsg []string

		err := c.ShouldBind(&payload)
		if err != nil {
			errMsg = append(errMsg, err.Error())
			c.HTML(200, "signup", gin.H{"messages": errMsg})
		} else {
			if payload.ID[:4] != "1420" {
				errMsg = append(errMsg, "Your ID does not belong to CSE Department")
			}
			if payload.Password != payload.ConfirmPassword {
				errMsg = append(errMsg, "Password Not match")
			}
			match1, _ := regexp.MatchString("[a-z]{1,3}", payload.Password)
			match2, _ := regexp.MatchString("[0-9]{1,3}", payload.Password)
			match3, _ := regexp.MatchString("[!@$&]{1,3}", payload.Password)
			if !match1 || !match2 || !match3 {
				errMsg = append(errMsg, "Password Not strong enough")
			}
			if len(errMsg) > 0 {
				c.HTML(200, "signup", gin.H{"messages": errMsg})
			} else {
				var newUser user.User

				hashedPass, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
				if err != nil {
					log.Println(err.Error())
				}

				authToken, err := bcrypt.GenerateFromPassword([]byte(payload.Username), bcrypt.DefaultCost)
				if err != nil {
					log.Println(err.Error())
				}

				newUser.FirstName = payload.FirstName
				newUser.LastName = payload.LastName
				newUser.Username = payload.Username
				newUser.Password = string(hashedPass)
				newUser.IsAuth = false
				newUser.DateCreated = time.Now()
				newUser.Email = payload.ID + "@iubat.edu"
				newUser.AuthToken = string(authToken)
				newUser.UserType = payload.UserType
				if payload.UserType == "judge" {
					newUser.IsJudge = true
					newUser.AllowJudge = false
				}

				Db := db.InitDb()
				Db.Create(&newUser)
				if err != nil {
					log.Println(err.Error())
					return
				}
				msg := "From: " + "IUBAT Beta" + "\n" +
					"To: " + payload.ID + "@iubat.edu" + "\n" +
					"Subject: Confirm Signup" + "\n" +
					"Please confirm your signup by following this link" + "\n" +
					"http://localhost:50051/auth/confirm/" + newUser.AuthToken + "\n" +
					"Team IUBAT Beta"

				go sendEmail(payload.ID+"@iubat.edu", msg)
				c.HTML(200, "signup", gin.H{"success": "true"})
			}
		}
	}
}

func sendEmail(addr, msg string) {
	auth := smtp.PlainAuth("", "cyantarek.cg@gmail.com", "zxcv5555", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "cyantarek.cg@gmail.com", []string{addr}, []byte(msg))
	if err != nil {
		log.Println(err.Error())
	}
}

func confirmationHandler(c *gin.Context) {
	authToken := c.Param("authToken")
	Db := db.InitDb()
	var count int
	Db.Model(user.User{}).Where("auth_token = ?", authToken).Count(&count)
	fmt.Println(count)

	if count == 1 {
		var checkUser user.User
		Db.Model(user.User{}).Where("auth_token = ?", authToken).First(&checkUser)

		if checkUser.IsAuth {
			c.HTML(200, "confirmation", gin.H{})
		} else {
			checkUser.IsAuth = true
			Db.Save(&checkUser)
			c.HTML(200, "confirmation", gin.H{"success": "true"})
		}
	} else {
		c.HTML(200, "confirmation", gin.H{})
	}
}

func recoveryHandler(c *gin.Context) {
	c.HTML(200, "recovery", nil)
}
