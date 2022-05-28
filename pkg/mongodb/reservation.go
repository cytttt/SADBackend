package mongodb

import (
	"SADBackend/model"
	"SADBackend/repo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoReservationRepo struct {
	collection *mongo.Collection
}

func newReservationRepository(db *mongo.Database) repo.ReservationRepo {
	return &mongoReservationRepo{collection: db.Collection("reservation")}
}

func (m *mongoReservationRepo) Exist(userID string, startAt time.Time, result interface{}) error {
	err := m.collection.FindOne(context.Background(), bson.M{"user_id": userID, "start_at": startAt}).Decode(result)

	return err
}

func (m *mongoReservationRepo) GetReservation(userID string, results interface{}) error {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))

	matchStage := bson.M{
		"$match": bson.M{
			"user_id":  bson.M{"$eq": userID},
			"start_at": bson.M{"$gte": time.Now().In(loc)},
		},
	}
	sortStage := bson.M{
		"$sort": bson.M{
			"start_at": 1,
		},
	}
	lookupStage1 := bson.M{
		"$lookup": bson.M{
			"from":         "machine",
			"localField":   "machine_id",
			"foreignField": "machine_id",
			"as":           "machines",
		},
	}
	lookupStage2 := bson.M{
		"$lookup": bson.M{
			"from":         "gym",
			"localField":   "machines.0.gym_id",
			"foreignField": "branch_gym_id",
			"as":           "gyms",
		},
	}
	pip := []bson.M{matchStage, sortStage, lookupStage1, lookupStage2}

	cursor, err := m.collection.Aggregate(context.Background(), pip)
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}

func (m *mongoReservationRepo) MakeReservation(userID, machineID string, start time.Time) error {
	newReservation := model.Reservation{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		MachineID: machineID,
		StartAt:   start,
	}
	if _, err := m.collection.InsertOne(context.Background(), newReservation); err != nil {
		return err
	}
	return nil
}

func (m *mongoReservationRepo) QueryExistReservation(machineIDs []string, lb, ub time.Time, results interface{}) error {
	filter := bson.M{
		"machine_id": bson.M{
			"$in": machineIDs,
		},
		"start_at": bson.M{"$gte": lb, "$lt": ub},
	}
	cursor, err := m.collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	if err := cursor.All(context.TODO(), results); err != nil {
		return err
	}
	return nil
}
