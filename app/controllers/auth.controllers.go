package controllers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"time"
	"fmt"
	"regexp"
	"strings"
	"net/smtp"
	"github.com/gin-gonic/gin"

	"../models"
	"../api"
	"log"
)

//Dependency Injection Pattern (injecting Store, SessionStore to Handler at runtime)
//Program to Interface (Store, SessionStore)
//Dependency De-Coupling
func (h *Handler) UserLogin(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		//api request
		if h.Type == "api" {
			type user struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			var payload user
			c.BindJSON(&payload)

			var checkUser models.User
			h.Store.GetSingle(&checkUser, "username = ?", payload.Username)

			if checkUser.ID != 0 {
				err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(payload.Password))
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
					c.Abort()
				} else {
					tok, err := api.GenerateToken(checkUser.Username, checkUser.UserType, 1, h.SignKey)
					if err != nil {
						c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
						c.Abort()
					}
					c.JSON(200, gin.H{"token": tok})
				}
			}
		//template request
		} else if h.Type == "template" {
			var userCheck models.User
			//
			fmt.Println(c.PostForm("username"))
			h.Store.GetSingle(&userCheck, "username = ?", c.PostForm("username"))
			if userCheck.Username == "" {
				c.HTML(200, "login", gin.H{})
				return
			} else if !userCheck.IsAuth {
				c.HTML(200, "login", gin.H{})
				return
			} else {
				start := time.Now()
				err := bcrypt.CompareHashAndPassword([]byte(userCheck.Password), []byte(c.PostForm("password")))
				//_, err := strconv.Atoi("a")
				fmt.Println(time.Since(start).Seconds())
				if err != nil {
					c.HTML(200, "login", gin.H{})
					return
				} else {
					userCheck.LastLogin = time.Now()
					userCheck.IsLoggedIn = true
					h.Store.UpdateData(&userCheck)

					h.SessionStore.SetSessionData(c, "user", userCheck)

					c.Redirect(http.StatusFound, "/")
				}
			}
		}


	} else if c.Request.Method == "GET" {
		c.HTML(200, "login", gin.H{"success": "true"})
	}
}

func (h *Handler) UserLogout(c *gin.Context) {
	h.SessionStore.SetSessionData(c, "user", nil)
	c.Redirect(http.StatusFound, "/auth/login")
}

func (h *Handler) UserSignup(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(200, "signup", gin.H{"messages": ""})
	} else if c.Request.Method == http.MethodPost {
		var errMsg []string

		if c.PostForm("id")[4:5] != "3" {
			errMsg = append(errMsg, "Your ID does not belong to CSE Department")
		}
		if c.PostForm("password") != c.PostForm("confirm_password") {
			errMsg = append(errMsg, "Password Not match")
		}
		match1, _ := regexp.MatchString("[a-z]{1,3}", c.PostForm("password"))
		match2, _ := regexp.MatchString("[0-9]{1,3}", c.PostForm("password"))
		match3, _ := regexp.MatchString("[!@$&]{1,3}", c.PostForm("password"))
		if !match1 || !match2 || !match3 {
			errMsg = append(errMsg, "Password Not strong enough")
		}
		if len(errMsg) > 0 {
			c.HTML(200, "signup", gin.H{"messages": errMsg})
		} else {
			var newUser models.User

			hashedPass, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err.Error())
			}

			authToken, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("username")), bcrypt.DefaultCost)
			if err != nil {
				log.Println(err.Error())
			}
			processedAuthToken := strings.Replace(string(authToken), "/", "", -1)

			newUser.FirstName = c.PostForm("first_name")
			newUser.LastName = c.PostForm("last_name")
			newUser.Username = c.PostForm("username")
			newUser.Password = string(hashedPass)
			newUser.IsAuth = false
			newUser.DateCreated = time.Now()
			newUser.Email = c.PostForm("id") + "@iubat.edu"
			newUser.AuthToken = processedAuthToken
			newUser.UserType = c.PostForm("user_type")
			newUser.LastLogin = time.Now()
			if newUser.UserType == "judge" {
				newUser.IsJudge = true
				newUser.AllowJudge = false
			} else if newUser.UserType == "participant" {
				newUser.IsParticipant = true
			}

			h.Store.InsertData(&newUser)
			if err != nil {
				log.Println(err.Error())
				return
			}
			msg := "From: " + "IUBAT Beta" + "\n" +
				"To: " + c.PostForm("id") + "@iubat.edu" + "\n" +
				"Subject: Confirm Signup" + "\n" +
				"Please confirm your signup by following this link" + "\n" +
				"http://localhost:50051/auth/confirm/" + newUser.AuthToken + "\n" +
				"Team IUBAT Beta"

			go sendEmail(c.PostForm("id")+"@iubat.edu", msg)
			c.HTML(200, "signup", gin.H{"success": "true"})
		}
	}
}

func sendEmail(addr, msg string) {
	auth := smtp.PlainAuth("", "cyantarek.cg@gmail.com", "zxcv5555", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "cyantarek.cg@gmail.com", []string{"cyantarek.cg@gmail.com"}, []byte(msg))
	if err != nil {
		fmt.Println(err.Error())
		log.Println(err.Error())
	}
	fmt.Println("email sent")
}

func (h *Handler) ConfirmationManager(c *gin.Context) {
	authToken := c.Param("authToken")

	var count int
	count = h.Store.GetCount("select count(*) from users where auth_token = ?", authToken)

	if count == 1 {
		var checkUser models.User
		h.Store.GetSingle(&checkUser, "select * from users where auth_token = ?", authToken)

		if checkUser.IsAuth {
			c.HTML(200, "confirmation", gin.H{})
		} else {
			checkUser.IsAuth = true
			h.Store.UpdateData(&checkUser)
			c.HTML(200, "confirmation", gin.H{"success": "true"})
		}
	} else {
		c.HTML(200, "confirmation", gin.H{})
	}
}

func (h *Handler) AccountRecovery(c *gin.Context) {
	c.HTML(200, "recovery", nil)
}
