package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type UpdateMachineStatusReq struct {
	MachineID string `json:"machine_id" example:"machine-1000001"`
	Amount    int    `json:"amount" example:"1"`
}

type MachineStatusWrapperResp struct {
	Category model.PartCategory  `json:"category"`
	Machines []MachineStatusResp `json:"machines"`
}
type MachineStatusResp struct {
	MachineID  string             `json:"machine_id"`
	Name       string             `json:"name"`
	WaitingPPL int                `json:"waiting_ppl"`
	Category   model.PartCategory `json:"category"`
}

// @Summary Get Machine Status Given Gym ID
// @Produce json
// @Tags All
// @param gym_id query string true "Gym id e.g. branch-1000001"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/gym/machine [get]
func GetMachineList(c *gin.Context) {
	gymID := c.Query("gym_id")
	err, results := findAllMachines(gymID)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
}

// @Summary Update Machine Status
// @Produce json
// @Tags Staff
// @Param UpdateMachineStatus body UpdateMachineStatusReq true "Machine_id, amount(int)"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/machine/status [put]
func UpdateMachineStatus(c *gin.Context) {
	var updateReq UpdateMachineStatusReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
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
	// get all machine in the same gym
	err, results := findAllMachines(machine.Gym)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
}

func findAllMachines(gymID string) (error, []MachineStatusResp) {
	cursor, err := mongodb.MachineCollection.Find(context.Background(), bson.M{"gym_id": gymID})
	if err != nil {
		return err, nil
	}
	var machines []model.Machine
	if err := cursor.All(context.TODO(), &machines); err != nil {
		return err, nil
	}

	var results []MachineStatusResp
	for _, machine := range machines {
		results = append(results, MachineStatusResp{
			MachineID:  machine.MachineID,
			Name:       machine.Name,
			Category:   machine.Category,
			WaitingPPL: machine.WaitingPPL,
		})
	}
	return nil, results
}

// @Summary Get Machine Status Given Gym ID Sorted by Category
// @Produce json
// @Tags All
// @Param gym_id path string true "Gym id e.g. branch-1000001"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/gym/machine/category/{gym_id} [get]
func GetMachineListByCategory(c *gin.Context) {
	log.Println("m")
	gymID := c.Param("gym_id")
	err, results := findAllMachines(gymID)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	categoryMapping := [...]model.PartCategory{model.PART_BACK, model.PART_CHEST, model.PART_CARDIO, model.PART_ABS, model.PART_LEG, model.PART_ARM, model.PART_HIPS}

	sortedResults := make(map[model.PartCategory][]MachineStatusResp)
	for _, m := range results {
		if m.Category == model.PART_ABS {
			sortedResults[model.PART_ABS] = append(sortedResults[model.PART_ABS], m)
		} else if m.Category == model.PART_ARM {
			sortedResults[model.PART_ARM] = append(sortedResults[model.PART_ARM], m)
		} else if m.Category == model.PART_BACK {
			sortedResults[model.PART_BACK] = append(sortedResults[model.PART_BACK], m)
		} else if m.Category == model.PART_CARDIO {
			sortedResults[model.PART_CARDIO] = append(sortedResults[model.PART_CARDIO], m)
		} else if m.Category == model.PART_CHEST {
			sortedResults[model.PART_CHEST] = append(sortedResults[model.PART_CHEST], m)
		} else if m.Category == model.PART_LEG {
			sortedResults[model.PART_LEG] = append(sortedResults[model.PART_LEG], m)
		} else if m.Category == model.PART_HIPS {
			sortedResults[model.PART_HIPS] = append(sortedResults[model.PART_HIPS], m)
		}
	}
	var finalResults []MachineStatusWrapperResp
	for _, category := range categoryMapping {
		finalResults = append(finalResults, MachineStatusWrapperResp{
			Category: category,
			Machines: sortedResults[category],
		})
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, finalResults)
}
