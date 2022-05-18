package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Backend Test *Frontend plz don't try
// @Produce json
// @Tags Backend
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/test [get]
func TTTTT(c *gin.Context) {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	// tmp := model.Reservation{
	// 	ID:          primitive.NewObjectID(),
	// 	UserID:      "meowmeow123",
	// 	MachineID:   "machine-1000001",
	// 	Category:    model.PART_Cardio,
	// 	MachineName: "treadmill-101",
	// 	StartAt:     time.Date(2023, time.Month(12), 10, 0, 0, 0, 0, loc),
	// }
	// _, err := mongodb.ReservationCollection.InsertOne(context.Background(), tmp)
	// if err != nil {
	// 	constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
	// 	return
	// }
	// filter := bson.M{"user_id": "meowmeow123"}
	// update := bson.M{"$set": bson.M{"statistics": model.Stat{
	// 	StayTime:   232,
	// 	Calories:   32424,
	// 	MostTrain:  model.PART_Cardio,
	// 	LeastTrain: model.PART_Arm,
	// }}}
	// _, err := mongodb.ClientCollection.UpdateOne(context.Background(), filter, update)
	// if err != nil {
	// 	constant.ResponseWithData(c, http.StatusOK, constant.ERROR, nil)
	// 	return
	// }
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
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
			return
		}
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}
