package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Backend Test *Frontend plz don't try
// @Produce json
// @Tags Backend
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/test [get]
func TTTTT(c *gin.Context) {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	tmp := model.Reservation{
		ID:          primitive.NewObjectID(),
		UserID:      "meowmeow123",
		MachineID:   "machine-1000001",
		Category:    model.PART_Cardio,
		MachineName: "treadmill-101",
		StartAt:     time.Date(2023, time.Month(12), 10, 0, 0, 0, 0, loc),
	}
	_, err := mongodb.ReservationCollection.InsertOne(context.Background(), tmp)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}
