package mongodb

import (
	"SADBackend/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGymRepo struct {
	collection *mongo.Collection
}

func newGymRepository(db *mongo.Database) repo.GymRepo {
	return &mongoGymRepo{collection: db.Collection("gym")}
}
