package contest

import "time"

type Contest struct {
	ID          int
	Name        string    `form:"contest-name"`
	ContestDate time.Time
	RunningTime int       `form:"running-time"`
}
