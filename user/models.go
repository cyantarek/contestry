package user

import "time"

type User struct {
	FirstName   string
	LastName    string
	Username    string
	ID          int
	Password    string
	IsAuth      bool
	DateCreated time.Time
	Email       string
	AuthToken   string
	IsAdmin     bool
	UserType    string
	IsJudge     bool
	AllowJudge  bool
	LastLogin   time.Time
	IsLoggedIn  bool
}
