package mongodb

import (
	"SADBackend/repo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoStaffRepo struct {
	collection *mongo.Collection
}

func newStaffRepository(db *mongo.Database) repo.StaffRepo {
	return &mongoStaffRepo{collection: db.Collection("staff")}
}

func (m *mongoStaffRepo) Exist(userID string, result interface{}) error {
	err := m.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(result)

	return err
}
