package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Machine struct {
	ID           primitive.ObjectID `bson:"_id"`
	MachineID    string             `bson:"user_id"`
	Name         string             `bson:"name"`
	Category     string             `bson:"category"`
	Gym          primitive.ObjectID `bson:"gym_id"`
	WaitingCount int                `bson:"waiting_count"`
	BoughtAt     time.Time          `bson:"bought_at"`
	MaintainedAt time.Time          `bson:"maintained_at"`
}
