package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	MachineID string             `bson:"machine_id"`
	StartAt   time.Time          `bson:"start_at"`
	Expired   bool               `bson:"expired"`
}

type ReservationResp struct {
	Category    string    `json:"category"`
	MachineID   string    `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	GymID       string    `json:"gym_id"`
	GymName     string    `json:"gym_name"`
	Date        time.Time `json:"date"`
}

type AggrReservationRes struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	MachineID   string             `bson:"machine_id"`
	Category    PartCategory       `bson:"category"`
	MachineName string             `bson:"machine_name"`
	StartAt     time.Time          `bson:"start_at"`
	Expired     bool               `bson:"expired"`
	Gyms        []BranchGym        `bson:"gyms"`
	Machines    []Machine          `bson:"machines"`
}

type GetAvailableTimeResp struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	MachineID string `json:"machine_id"`
}
