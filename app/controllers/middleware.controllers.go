package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

func (h *Handler) JwtValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		//first check if the token is valid or invalid by querying jwt whitelist redis db
		//if valid, then go on next
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
					h.ContextData = make(map[string]interface{})
					h.ContextData["user_id"] = tok.Claims.(jwt.MapClaims)["user_id"]
					h.ContextData["username"] = tok.Claims.(jwt.MapClaims)["username"]
					h.ContextData["user_type"] = tok.Claims.(jwt.MapClaims)["user_type"]
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

func (h *Handler) AdminValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.ContextData["user_type"] != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You're not allowed to do this"})
			c.Abort()
			return
		}
	}
}

func (h *Handler) JudgeValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.ContextData["user_type"] != "judge" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You're not allowed to do this"})
			c.Abort()
			return
		}
	}
}

func (h *Handler) ParticipantValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.ContextData["user_type"] != "participant" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You're not allowed to do this"})
			c.Abort()
			return
		}
	}
}
