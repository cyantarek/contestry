package session_mgr

import (
	"github.com/gin-gonic/gin"
	"../models"
)

//wrapping concrete dependencies with application specific struct
type MockSession struct {}

func NewMockSession(name string, secret string, r *gin.Engine) SessionStore {
	return &MockSession{}
}

func (m *MockSession) GetSessionData(ctx *gin.Context, key string) interface{} {
	var user models.User
	user.Username = "testUser"
	user.ID = 001
	user.IsParticipant = true
	user.IsLoggedIn = true
	user.UserType = "participant"
	user.Email = "testUser@gmail.com"
	user.FirstName = "Test"
	user.LastName = "User"
	user.IsAdmin = true

	return user
}

func (m *MockSession) SetSessionData(ctx *gin.Context, key string, data interface{}) {

}


