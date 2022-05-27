package repo

import (
	"SADBackend/model"

	"go.mongodb.org/mongo-driver/bson"
)

type ClientRepo interface {
	Exist(userID string, result interface{}) error
	UpdateClientInfo(userID string, update bson.M, result interface{}) error
	Signup(newClient model.Client) error
}

type StaffRepo interface {
}

type ReservationRepo interface {
}

type GymRepo interface {
}

type MachineRepo interface {
}

type AttendanceRepo interface {
}

var Client ClientRepo
var Staff StaffRepo
var Reservation ReservationRepo
var Gym GymRepo
var Machine MachineRepo
var Attendance AttendanceRepo
