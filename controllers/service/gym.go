package service

import "SADBackend/model"

func PostprocessGymList(raw []model.BranchGym) ([]model.GetGymListResp, error) {
	var results []model.GetGymListResp
	for _, gym := range raw {
		status := "uncrowded"
		if float64(gym.CurrentNumberPeople) > float64(gym.Info.ClientNumberLimit)*0.8 {
			status = "crowded"
		}
		results = append(results, model.GetGymListResp{
			BranchGymID:      gym.BranchGymID,
			Name:             gym.Name,
			Address:          gym.Address,
			Status:           status,
			AvailableMachine: gym.AvailableMachine,
		})
	}
	return results, nil
}
