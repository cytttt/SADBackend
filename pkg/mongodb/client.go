package mongodb

import (
	"SADBackend/model"
	"SADBackend/repo"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClientRepo struct {
	collection *mongo.Collection
}

func newClientRepository(db *mongo.Database) repo.ClientRepo {
	return &mongoClientRepo{collection: db.Collection("client")}
}

// check client exist or not
func (m *mongoClientRepo) Exist(userID string, result interface{}) error {
	err := m.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(result)

	return err
}

// signup
func (m *mongoClientRepo) Signup(newClient model.Client) error {
	if _, err := m.collection.InsertOne(context.Background(), newClient); err != nil {
		log.Printf("User ID: %s / Error: %s", newClient.UserID, err)
		return err
	}
	return nil
}

//
func (m *mongoClientRepo) UpdateClientInfo(userID string, update bson.M, result interface{}) error {
	opt := options.FindOneAndUpdate()
	opt.SetUpsert(true)
	opt.SetReturnDocument(options.After)
	filter := bson.M{"user_id": userID}
	err := m.collection.FindOneAndUpdate(context.Background(), filter, update, opt).Decode(result)
	return err
}

func (m *mongoClientRepo) DeleteClient(userID string) error {
	_, err := m.collection.DeleteOne(context.Background(), bson.M{"user_id": userID})
	return err
}
