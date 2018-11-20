package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"../models"
	"github.com/mitchellh/mapstructure"
	"fmt"
)

//Adapter or Wrapper or Decorator Pattern
func (h *Handler) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser := h.SessionStore.GetSessionData(c, "user")
		if loginUser == nil {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
		} else {
			h.ContextData = loginUser
			fmt.Println(loginUser)
			c.Next()
		}
	}
}

func (h *Handler) UnAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		loginUser := h.SessionStore.GetSessionData(c, "user")
		if loginUser != nil {
			c.Redirect(http.StatusFound, "/")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func (h *Handler) TypeChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.URL.Path) > 1 && c.Request.URL.Path[:4] == "/api" {
			h.Type = "api"
			c.Next()
		} else {
			h.Type = "template"
			c.Next()
		}
	}
}

func (h *Handler) JwtValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader != "" {
			bearerTok := strings.Split(authHeader, " ")

			if len(bearerTok) == 2 {
				tok, err := jwt.Parse(bearerTok[1], func(token *jwt.Token) (interface{}, error) {
					return h.VerifyKey, nil
				})

				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					c.Abort()
					return
				}

				if tok.Valid {
					var user models.User
					claims := tok.Claims
					mapstructure.Decode(claims, &user)
					h.ContextData = user
					c.Next()
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Auth Token"})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token Required"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Un-Authorized"})
			c.Abort()
			return
		}
	}
}
