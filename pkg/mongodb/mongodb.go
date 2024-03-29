package mongodb

import (
	"SADBackend/repo"
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

	// refactor v1
	// repo.Client = newClientRepository(DB)
	// repo.Staff = newStaffRepository(DB)
	// repo.Gym = newGymRepository(DB)
	// repo.Machine = newMachineRepository(DB)
	// repo.Reservation = newReservationRepository(DB)
	// repo.Attendance = newAttendanceRepository(DB)

	// refactor v2
	repo.RepoInstance.Client = newClientRepository(DB)
	repo.RepoInstance.Staff = newStaffRepository(DB)
	repo.RepoInstance.Gym = newGymRepository(DB)
	repo.RepoInstance.Machine = newMachineRepository(DB)
	repo.RepoInstance.Reservation = newReservationRepository(DB)
	repo.RepoInstance.Attendance = newAttendanceRepository(DB)

	log.Printf("[info] MongoDB initialization is done")
}

func Dispose() {
	log.Println("shut down MongoDB connection")
	if err := MongoInstance.Client.Disconnect(context.TODO()); err != nil {
		log.Printf("mongo disconnect error: %v", err)
	}
}
