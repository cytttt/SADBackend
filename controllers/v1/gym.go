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

type GetGymListResp struct {
	BranchGymID string `json:"branch_gym_id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Status      string `json:"status"`
}

// @Summary Get Gym List
// @Produce json
// @Tags All
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/gym/list [get]
func GetGymList(c *gin.Context) {
	cursor, err := mongodb.GymCollection.Find(context.Background(), bson.M{})
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, nil)
		return
	}
	var gyms []model.BranchGym
	if err := cursor.All(context.TODO(), &gyms); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, nil)
		return
	}
	var results []GetGymListResp
	for _, gym := range gyms {
		status := "uncrowded"
		if float64(gym.CurrentNumberPeople) > float64(gym.Info.ClientNumberLimit)*0.8 {
			status = "crowded"
		}
		results = append(results, GetGymListResp{
			BranchGymID: gym.BranchGymID,
			Name:        gym.Name,
			Address:     gym.Address,
			Status:      status,
		})
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
}
