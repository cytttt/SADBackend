package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateMachineStatusReq struct {
	MachineID string `json:"machine_id"`
	Amount    int    `json:"amount"`
}

type UpdateMachineStatusResp struct {
	MachineID  string `json:"machine_id"`
	Name       string `bson:"name"`
	WaitingPPL int    `json:"waiting_ppl"`
	Category   string `bson:"category"`
}

// @Summary Update machine status
// @Tags Staff
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/machine/status [put]
func UpdateMachineStatus(c *gin.Context) {
	var updateReq UpdateMachineStatusReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, nil)
		return
	}
	// update
	var machine model.Machine
	update := bson.M{"$inc": bson.M{"waiting_ppl": updateReq.Amount}}
	err := mongodb.MachineCollection.FindOneAndUpdate(context.Background(), bson.M{"machine_id": updateReq.MachineID}, update).Decode(&machine)
	if err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	// get all machine
	cursor, err := mongodb.MachineCollection.Find(context.Background(), bson.M{"gym_id": machine.Gym})
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, nil)
		return
	}
	var machines []model.Machine
	if err := cursor.All(context.TODO(), &machines); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, nil)
		return
	}
	var results []UpdateMachineStatusResp
	for _, machine := range machines {
		results = append(results, UpdateMachineStatusResp{
			MachineID:  machine.MachineID,
			Name:       machine.Name,
			Category:   machine.Category,
			WaitingPPL: machine.WaitingPPL,
		})
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
}
