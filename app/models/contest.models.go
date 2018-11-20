package models

import (
	"time"
)

type Contest struct {
	ID            int
	Name          string
	ContestDate   time.Time
	RunningTime   int
	Slug          string
	Status        string
	QuestionLimit int
	Type          string
	Notes         string
}

type Question struct {
	ID           int
	Title        string
	Description  string
	InputFormat  string
	OutputFormat string
	Constrains   string
	SampleInput  string
	SampleOutput string
	Input        string
	CorrectAns   string
	Point        float32
	ContestID    int
	UserID       int
}

type Solution struct {
	ID         int
	Result     string
	ContestID  string
	QuestionID string
	UserID     string
	ExecTime   string
	Point      float32
}
