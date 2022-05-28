package mongodb

import (
	"SADBackend/repo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoMachineRepo struct {
	collection *mongo.Collection
}

func newMachineRepository(db *mongo.Database) repo.MachineRepo {
	return &mongoMachineRepo{collection: db.Collection("machine")}
}

func (m *mongoMachineRepo) MachineList(gymID string, results interface{}) error {
	cursor, err := m.collection.Find(context.Background(), bson.M{"gym_id": gymID, "reservation_only": false})
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}

func (m *mongoMachineRepo) UpdateAmount(machineID string, amount int, result interface{}) error {
	update := bson.M{"$inc": bson.M{"waiting_ppl": amount}}
	err := m.collection.FindOneAndUpdate(context.Background(), bson.M{"machine_id": machineID}, update).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoMachineRepo) GetAvailableMachines(gymID, machineName string, results interface{}) error {
	machineFilter := bson.M{
		"reservation_only": true,
		"gym_id":           gymID,
		"name": primitive.Regex{
			Pattern: "^" + machineName,
			Options: "",
		},
	}
	cursor, err := m.collection.Find(context.Background(), machineFilter)
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}
