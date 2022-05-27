package mongodb

import (
	"SADBackend/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoMachineRepo struct {
	collection *mongo.Collection
}

func newMachineRepository(db *mongo.Database) repo.MachineRepo {
	return &mongoClientRepo{collection: db.Collection("machine")}
}
