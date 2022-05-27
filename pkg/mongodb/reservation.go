package mongodb

import (
	"SADBackend/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoReservationRepo struct {
	collection *mongo.Collection
}

func newReservationRepository(db *mongo.Database) repo.ReservationRepo {
	return &mongoClientRepo{collection: db.Collection("reservation")}
}
