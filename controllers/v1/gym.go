package v1

import (
	"SADBackend/constant"
	"SADBackend/controllers/service"
	"SADBackend/model"
	"SADBackend/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Gym List
// @Produce json
// @Tags All
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/gym/list [get]
func GetGymList(c *gin.Context, gymDB repo.GymRepo) {
	var gyms []model.BranchGym
	if err := gymDB.GymList(&gyms); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	results, _ := service.PostprocessGymList(gyms)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
}
