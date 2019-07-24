package controllers

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"../models"
	"../helper"
	"regexp"
	"strings"
	"time"
	"log"
	"strconv"
)

func (h *Handler) ApiLoginUser(c *gin.Context) {
	type payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var payloadData payload
	err := c.BindJSON(&payloadData)
	if err != nil {
		log.Println(err.Error())
		return
	}

	var checkUser models.User
	h.Store.GetSingle(&checkUser, "username = ?", payloadData.Username)

	if checkUser.ID != 0 {
		if checkUser.IsAuth != true {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated, please confirm your email"})
			c.Abort()
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(payloadData.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			c.Abort()
		} else {
			tok, err := helper.GenerateToken(checkUser.Username, checkUser.UserType, checkUser.ID, 1, h.SignKey)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				c.Abort()
			}
			//c.SetCookie("token", tok, 120, "/", "127.0.0.1", false, false)
			//save the token to jwt token whitelist to make it valid (or logged in)
			c.JSON(200, gin.H{"token": tok})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
}

func (h *Handler) ApiSignupUser(c *gin.Context) {
	type payload struct {
		FirstName       string `json:"first_name" binding:"required"`
		LastName        string `json:"last_name" binding:"required"`
		Username        string `json:"username" binding:"required"`
		UniID           string `json:"uni_id" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		Email           string `json:"email" binding:"required"`
		UserType        string `json:"user_type" binding:"required"`
	}

	var payloadData payload
	err := c.BindJSON(&payloadData)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "All fields are not filled"})
		return
	}

	checkUser := h.Store.GetCount("select count(*) from users where username = ? or uni_id = ? or email = ?", payloadData.Username, payloadData.UniID, payloadData.Email)
	if checkUser > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already exists"})
		return
	}

	if payloadData.UniID[4:5] != "3" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Your ID does not belong to CSE Department"})
	}
	if payloadData.Password != payloadData.ConfirmPassword {
		c.JSON(http.StatusForbidden, gin.H{"error": "Password Not match"})
	}
	match1, _ := regexp.MatchString("[a-z]{1,3}", payloadData.Password)
	match2, _ := regexp.MatchString("[0-9]{1,3}", payloadData.Password)
	match3, _ := regexp.MatchString("[!@$&]{1,3}", payloadData.Password)
	if !match1 || !match2 || !match3 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Password Not strong enough"})
	}

	var newUser models.User

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(payloadData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
	}

	authToken, err := bcrypt.GenerateFromPassword([]byte(payloadData.Username), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return
	}
	processedAuthToken := strings.Replace(string(authToken), "/", "", -1)

	newUser.FirstName = payloadData.FirstName
	newUser.LastName = payloadData.LastName
	newUser.Username = payloadData.Username
	newUser.Password = string(hashedPass)
	newUser.IsAuth = false
	newUser.DateJoined = time.Now()
	newUser.Email = payloadData.UniID + "@iubat.edu"
	newUser.AuthToken = processedAuthToken
	newUser.UserType = payloadData.UserType
	newUser.LastLogin = time.Now()
	univId, _ := strconv.Atoi(payloadData.UniID)
	newUser.UniID = univId

	h.Store.InsertData(&newUser)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//msg := "From: " + "IUBAT Beta" + "\n" +
	//	"To: " + c.PostForm("id") + "@iubat.edu" + "\n" +
	//	"Subject: Confirm Signup" + "\n" +
	//	"Please confirm your signup by following this link" + "\n" +
	//	"http://localhost:50051/auth/confirm/" + newUser.AuthToken + "\n" +
	//	"Team IUBAT Beta"
	//
	//go sendEmail(c.PostForm("id")+"@iubat.edu", msg)
	c.JSON(http.StatusCreated, gin.H{"msg": "new user created"})
}
