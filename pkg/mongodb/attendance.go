package mongodb

import (
	"SADBackend/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAttendanceRepo struct {
	collection *mongo.Collection
}

func newAttendanceRepository(db *mongo.Database) repo.AttendanceRepo {
	return &mongoAttendanceRepo{collection: db.Collection("attendance")}
}
