package service

import (
	"SADBackend/model"
	"SADBackend/repo"
	"log"
)

func PostprocessMachineList(raw []model.Machine) ([]model.MachineStatusResp, error) {
	var results []model.MachineStatusResp
	for _, machine := range raw {
		results = append(results, model.MachineStatusResp{
			MachineID:  machine.MachineID,
			Name:       machine.Name,
			Category:   machine.Category,
			WaitingPPL: machine.WaitingPPL,
		})
	}
	return results, nil
}

func PostprocessMachineCategory(raw []model.MachineStatusResp) ([]model.MachineStatusWrapperResp, error) {
	categoryMapping := [...]model.PartCategory{model.PART_BACK, model.PART_CHEST, model.PART_CARDIO, model.PART_ABS, model.PART_LEG, model.PART_ARM, model.PART_HIPS}

	sortedResults := make(map[model.PartCategory][]model.MachineStatusResp)
	for _, m := range raw {
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

	var results []model.MachineStatusWrapperResp
	for _, category := range categoryMapping {
		results = append(results, model.MachineStatusWrapperResp{
			Category: category,
			Machines: sortedResults[category],
		})
	}
	return results, nil
}

func FindAvailableMachine(gymID, machineName string, machineDB repo.MachineRepo) ([]string, error) {
	var machines []model.Machine
	if err := machineDB.GetAvailableMachines(gymID, machineName, &machines); err != nil {
		return []string{}, err
	}
	log.Println(machines)
	var results []string
	for _, i := range machines {
		results = append(results, i.MachineID)
	}
	return results, nil
}
