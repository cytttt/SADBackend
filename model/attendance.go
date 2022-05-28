package model

import "time"

type Attendance struct {
	Enter    time.Time `bson:"enter"`
	UserID   string    `bson:"user_id"`
	StayTime int       `bson:"stay_time"` // in seconds
}

type CompanyStatResp struct {
	Date        string  `json:"date"`
	Attendance  int     `json:"attendance_count"`
	AvgStayTime float32 `json:"avg_stay_hour"`
}

type StatInSecond struct {
	Date            string  `bson:"_id"`
	AttendanceCount int     `bson:"attendance_count"`
	AvgStaySecond   float64 `bson:"avg_stay_second"`
}
