package mongodb

import (
	"SADBackend/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStaffRepo struct {
	collection *mongo.Collection
}

func newStaffRepository(db *mongo.Database) repo.StaffRepo {
	return &mongoClientRepo{collection: db.Collection("staff")}
}
