package mongodb

import (
	"SADBackend/repo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAttendanceRepo struct {
	collection *mongo.Collection
}

func newAttendanceRepository(db *mongo.Database) repo.AttendanceRepo {
	return &mongoAttendanceRepo{collection: db.Collection("attendance")}
}

func (m *mongoAttendanceRepo) CompanyStat7days(results interface{}) error {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	cur := time.Now().In(loc)
	year, month, day := cur.Date()
	ub := time.Date(year, month, day, 0, 0, 0, 0, loc)
	lb := ub.AddDate(0, 0, -7)
	matchStage := bson.M{
		"$match": bson.M{
			"enter": bson.M{"$gte": lb, "$lt": ub},
		},
	}
	groupStage := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y/%m/%d",
					"date":   "$enter",
				}},
			"attendance_count": bson.M{"$sum": 1},
			"avg_stay_second":  bson.M{"$avg": "$stay_time"},
		},
	}
	pip := []bson.M{matchStage, groupStage}
	cursor, err := m.collection.Aggregate(context.Background(), pip)
	if err != nil {
		return err
	}

	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}
