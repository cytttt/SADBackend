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
}

type BranchInfo struct {
	CreatedAt         time.Time `bson:"created_at"`
	ClientNumberLimit int       `bson:"client_number_limit"`
}
