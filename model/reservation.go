package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	MachineID   string             `bson:"machine_id"`
	Category    PartCategory       `bson:"category"`
	MachineName string             `bson:"machine_name"`
	StartAt     time.Time          `bson:"start_at"`
	Expired     bool               `bson:"expired"`
}
