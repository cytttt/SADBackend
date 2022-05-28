package v1

import (
	"SADBackend/constant"
	"SADBackend/controllers/service"
	"SADBackend/model"
	"SADBackend/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateMachineStatusReq struct {
	MachineID string `json:"machine_id" example:"machine-1000001"`
	Amount    int    `json:"amount" example:"1"`
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
func GetMachineList(c *gin.Context, machineDB repo.MachineRepo) {
	gymID := c.Query("gym_id")

	var machines []model.Machine
	if err := machineDB.MachineList(gymID, &machines); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	res, _ := service.PostprocessMachineList(machines)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

// @Summary Update Machine Status
// @Produce json
// @Tags Staff
// @Param UpdateMachineStatus body UpdateMachineStatusReq true "Machine_id, amount(int)"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/machine/status [put]
func UpdateMachineStatus(c *gin.Context, machineDB repo.MachineRepo) {
	var updateReq UpdateMachineStatusReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	var machine *model.Machine
	err := machineDB.UpdateAmount(updateReq.MachineID, updateReq.Amount, &machine)
	if err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	var machines []model.Machine
	if err := machineDB.MachineList(machine.Gym, &machines); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	res, _ := service.PostprocessMachineList(machines)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

// @Summary Get Machine Status Given Gym ID Sorted by Category
// @Produce json
// @Tags All
// @Param gym_id path string true "Gym id e.g. branch-1000001"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/gym/machine/category/{gym_id} [get]
func GetMachineListByCategory(c *gin.Context, machineDB repo.MachineRepo) {
	gymID := c.Param("gym_id")

	var machines []model.Machine
	if err := machineDB.MachineList(gymID, &machines); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	rawRes, _ := service.PostprocessMachineList(machines)
	res, _ := service.PostprocessMachinecategory(rawRes)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}
