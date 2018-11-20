package models

import "time"

type User struct {
	FirstName     string
	LastName      string
	Username      string
	ID            int
	Password      string
	IsAuth        bool
	DateCreated   time.Time
	Email         string
	AuthToken     string
	IsAdmin       bool
	UserType      string
	IsJudge       bool
	IsParticipant bool
	AllowJudge    bool
	LastLogin     time.Time
	IsLoggedIn    bool
	Teams        []Team `gorm:"many2many:user_team;"`
}

type Team struct {
	ID       int
	Name     string
	TagLine  string
	Status   string
	LeaderID int
	Slug     string
	DateCreated   time.Time
}
