package mongodb

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoInstance *MongoAccess

type MongoAccess struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var ClientCollection *mongo.Collection
var StaffCollection *mongo.Collection
var ReservationCollection *mongo.Collection
var GymCollection *mongo.Collection
var MachineCollection *mongo.Collection
var AttendanceCollection *mongo.Collection

func Init() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(viper.GetString("MONGO_DB_CONNECTION")))
	if err != nil {
		panic(err)
	}

	DB := client.Database(viper.GetString("MONGO_DB_NAME"))
	MongoInstance = &MongoAccess{
		Client: client,
		DB:     DB,
	}

	// Setup Collection
	ClientCollection = DB.Collection("client")
	StaffCollection = DB.Collection("staff")
	ReservationCollection = DB.Collection("reservation")
	GymCollection = DB.Collection("gym")
	MachineCollection = DB.Collection("machine")
	AttendanceCollection = DB.Collection("attendance")

	log.Printf("[info] MongoDB initialization is done")
}

func Dispose() {
	log.Println("shut down MongoDB connection")
	if err := MongoInstance.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("mongo disconnect error: %v", err)
	}
}
