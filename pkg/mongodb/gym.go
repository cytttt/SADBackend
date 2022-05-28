package mongodb

import (
	"SADBackend/repo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGymRepo struct {
	collection *mongo.Collection
}

func newGymRepository(db *mongo.Database) repo.GymRepo {
	return &mongoGymRepo{collection: db.Collection("gym")}
}

func (m *mongoGymRepo) GymList(results interface{}) error {
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}
