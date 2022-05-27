package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"SADBackend/repo"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Backend Test *Frontend plz don't try
// @Produce json
// @Tags Backend
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/test [get]
func TTTTT(c *gin.Context) {
	// if err := addReservation(); err != nil {
	// 	constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
	// 	return
	// }
	// if err := addStuffStat(); err != nil {
	// 	constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
	// 	return
	// }
	// if err := addAvailableMachine(); err != nil {
	// 	constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
	// 	return
	// }
	err := repo.Client.Exist("meowmeow123", struct{}{})
	if err == nil {
		log.Println("nice")
	}
	var a *model.Client
	err = repo.Client.Exist("meowmeow123", &a)
	if err == nil {
		log.Println("nice2")
		log.Println(*a)
	}
	var b *ClientInfoResp
	err = repo.Client.Exist("meowmeow123", &b)
	if err == nil {
		log.Println("nice2")
		log.Println(*b)
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}

func addAvailableMachine() error {
	tmpAvailableMachine := model.MachineList{
		ABS:    []model.MachineName{model.MACHINE_AB_BENCH, model.MACHINE_AB_COASTER},
		ARM:    []model.MachineName{model.MACHINE_BICEP_VEST, model.MACHINE_ROWER},
		BACK:   []model.MachineName{model.MACHINE_PULLDOWN, model.MACHINE_POWER_TOWER},
		CARDIO: []model.MachineName{model.MACHINE_SPINNER, model.MACHINE_TREADMILL},
		CHEST:  []model.MachineName{model.MACHINE_CHEST_PRESS, model.MACHINE_PEC_DECK},
		HIPS:   []model.MachineName{model.MACHINE_DONKEY_KICK, model.MACHINE_GHD},
		LEG:    []model.MachineName{model.MACHINE_HACK_SQUAT, model.MACHINE_LEG_PRESS},
	}
	ml := [...]string{"branch-1000001", "branch-1000002"}
	opt := &options.UpdateOptions{}
	opt.SetUpsert(true)
	for _, i := range ml {
		update := bson.M{
			"$set": bson.M{
				"available_machine": tmpAvailableMachine,
			},
		}
		_, err := mongodb.GymCollection.UpdateOne(context.Background(), bson.M{"branch_gym_id": i}, update, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

func addReservation() error {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	tmp := model.Reservation{
		ID:        primitive.NewObjectID(),
		UserID:    "meowmeow123",
		MachineID: "machine-1000001",
		StartAt:   time.Date(2023, time.Month(12), 10, 0, 0, 0, 0, loc),
	}
	_, err := mongodb.ReservationCollection.InsertOne(context.Background(), tmp)
	if err != nil {
		return err
	}
	filter := bson.M{"user_id": "meowmeow123"}
	update := bson.M{"$set": bson.M{"statistics": model.Stat{
		StayTime:   232,
		Calories:   32424,
		MostTrain:  model.PART_CARDIO,
		LeastTrain: model.PART_ARM,
	}}}
	_, err = mongodb.ClientCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func addStuffStat() error {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	cases := [...]model.Attendance{
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -1),
			StayTime: 3600,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -2),
			StayTime: 4800,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -3),
			StayTime: 7200,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -4),
			StayTime: 5600,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -5),
			StayTime: 1800,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -6),
			StayTime: 3800,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -7),
			StayTime: 4400,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -2),
			StayTime: 4200,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -4),
			StayTime: 7200,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -6),
			StayTime: 7200,
		},
		{
			Enter:    time.Now().In(loc).AddDate(0, 0, -7),
			StayTime: 7200,
		},
	}
	for _, i := range cases {
		_, err := mongodb.AttendanceCollection.InsertOne(context.Background(), i)
		if err != nil {
			return err
		}
	}
	return nil

}
