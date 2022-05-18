package model

import "time"

type Attendance struct {
	Enter    time.Time `bson:"enter"`
	UserID   string    `bson:"user_id"`
	StayTime int       `bson:"stay_time"` // in seconds
}
