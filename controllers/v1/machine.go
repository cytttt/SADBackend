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
	categoryMapping := [...]model.PartCategory{model.PART_Back, model.PART_Chest, model.PART_Cardio, model.PART_ABS, model.PART_Leg, model.PART_Arm, model.PArt_Hips}

	sortedResults := make(map[model.PartCategory][]MachineStatusResp)
	for _, m := range results {
		if m.Category == model.PART_ABS {
			sortedResults[model.PART_ABS] = append(sortedResults[model.PART_ABS], m)
		} else if m.Category == model.PART_Arm {
			sortedResults[model.PART_Arm] = append(sortedResults[model.PART_Arm], m)
		} else if m.Category == model.PART_Back {
			sortedResults[model.PART_Back] = append(sortedResults[model.PART_Back], m)
		} else if m.Category == model.PART_Cardio {
			sortedResults[model.PART_Cardio] = append(sortedResults[model.PART_Cardio], m)
		} else if m.Category == model.PART_Chest {
			sortedResults[model.PART_Chest] = append(sortedResults[model.PART_Chest], m)
		} else if m.Category == model.PART_Leg {
			sortedResults[model.PART_Leg] = append(sortedResults[model.PART_Leg], m)
		} else if m.Category == model.PArt_Hips {
			sortedResults[model.PArt_Hips] = append(sortedResults[model.PArt_Hips], m)
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
