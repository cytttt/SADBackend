package repo

import (
	"SADBackend/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ClientRepo interface {
	Exist(userID string, result interface{}) error
	UpdateClientInfo(userID string, update bson.M, result interface{}) error
	Signup(newClient model.Client) error
}

type StaffRepo interface {
	Exist(userID string, result interface{}) error
}

type ReservationRepo interface {
	Exist(userID string, startAt time.Time, result interface{}) error
	GetReservation(userID string, results interface{}) error
	MakeReservation(userID, machineID string, start time.Time) error
	QueryExistReservation(machineIDs []string, lb, ub time.Time, results interface{}) error
}

type GymRepo interface {
	GymList(results interface{}) error
}

type MachineRepo interface {
	MachineList(gymID string, results interface{}) error
	UpdateAmount(machineID string, amount int, result interface{}) error
	GetAvailableMachines(gymID, machineName string, results interface{}) error
}

type AttendanceRepo interface {
	CompanyStat7days(results interface{}) error
}

var Client ClientRepo
var Staff StaffRepo
var Reservation ReservationRepo
var Gym GymRepo
var Machine MachineRepo
var Attendance AttendanceRepo
