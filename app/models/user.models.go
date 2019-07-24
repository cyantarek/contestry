package models

import "time"

type User struct {
	ID         int
	FirstName  string
	LastName   string
	Username   string
	UniID      int
	Password   string
	IsAuth     bool
	DateJoined time.Time
	Email      string
	AuthToken  string
	UserType   string
	AllowJudge bool
	LastLogin  time.Time
	//Teams        []Team `gorm:"many2many:user_team;"`
}

type Team struct {
	ID          int
	Name        string
	TagLine     string
	Status      string
	LeaderID    int
	Slug        string
	DateCreated time.Time
}
