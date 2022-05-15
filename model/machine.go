package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Machine struct {
	ID           primitive.ObjectID `bson:"_id"`
	MachineID    string             `bson:"machine_id"`
	Name         string             `bson:"name"`
	Category     PartCategory       `bson:"category"`
	Gym          string             `bson:"gym_id"`
	WaitingPPL   int                `bson:"waiting_ppl"`
	BoughtAt     time.Time          `bson:"bought_at"`
	MaintainedAt time.Time          `bson:"maintained_at"`
}

type PartCategory string

const (
	PART_Back   PartCategory = "client"
	PART_Chest  PartCategory = "staff"
	PART_Cardio PartCategory = "cardio"
	PART_ABS    PartCategory = "abs"
	PART_Leg    PartCategory = "leg"
	PART_Arm    PartCategory = "arm"
	PArt_Hips   PartCategory = "hips"
)
