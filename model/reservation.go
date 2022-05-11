package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type reservation struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	MachineID   string             `bson:"user_id"`
	Category    string             `bson:"category"`
	MachineName string             `bson:"machine_name"`
	StartAt     time.Time          `bson:"start_at"`
}
