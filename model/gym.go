package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BranchGym struct {
	ID                  primitive.ObjectID `bson:"_id"`
	BranchGymID         string             `bson:"branch_gym_id"`
	Name                string             `bson:"name"`
	Info                BranchInfo         `bson:"info"`
	Address             string             `bson:"address"`
	CurrentNumberPeople int                `bson:"current_number_people"`
	AvailableMachine    MachineList        `bson:"available_machine"`
}

type BranchInfo struct {
	CreatedAt         time.Time `bson:"created_at"`
	ClientNumberLimit int       `bson:"client_number_limit"`
}

type MachineList struct {
	ABS    []MachineName `bson:"abs" json:"abs"`
	ARM    []MachineName `bson:"arm" json:"arm"`
	BACK   []MachineName `bson:"back" json:"back"`
	CARDIO []MachineName `bson:"cardio" json:"cardio"`
	CHEST  []MachineName `bson:"chest" json:"chest"`
	HIPS   []MachineName `bson:"hips" json:"hips"`
	LEG    []MachineName `bson:"leg" json:"leg"`
}

type GetGymListResp struct {
	BranchGymID      string      `json:"branch_gym_id"`
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	Status           string      `json:"status"`
	AvailableMachine MachineList `json:"available_machine"`
}
