package service

import (
	"SADBackend/model"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeMachineRepo struct{}

func (_ FakeMachineRepo) MachineList(gymID string, results interface{}) error {
	return nil
}
func (_ FakeMachineRepo) UpdateAmount(machineID string, amount int, result interface{}) error {
	return nil
}
func (_ FakeMachineRepo) GetAvailableMachines(gymID, machineName string, results interface{}) error {
	switch results.(type) {
	case *[]model.Machine:
		res := results.(*[]model.Machine)
		*res = append(*res, model.Machine{
			MachineID: "test-machine-1",
		})
		*res = append(*res, model.Machine{
			MachineID: "test-machine-2",
		})
	default:
		log.Print("unknown type")
	}
	return nil
}

func TestPostprocessMachineCategory(t *testing.T) {
	tests := []struct {
		name    string
		args    []model.MachineStatusResp
		want    []model.MachineStatusWrapperResp
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: correct order",
			args: []model.MachineStatusResp{
				{
					MachineID: "test-1",
					Category:  model.PART_BACK,
				},
				{
					MachineID: "test-2",
					Category:  model.PART_CARDIO,
				},
				{
					MachineID: "test-3",
					Category:  model.PART_BACK,
				},
			},
			want: []model.MachineStatusWrapperResp{
				{
					Category: model.PART_BACK,
					Machines: []model.MachineStatusResp{
						{
							MachineID: "test-1",
							Category:  model.PART_BACK,
						},
						{
							MachineID: "test-3",
							Category:  model.PART_BACK,
						},
					},
				},
				{},
				{
					Category: model.PART_CARDIO,
					Machines: []model.MachineStatusResp{
						{
							MachineID: "test-2",
							Category:  model.PART_CARDIO,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostprocessMachineCategory(tt.args)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want[0].Machines[0], got[0].Machines[0])
			assert.Equal(t, tt.want[0].Machines[1], got[0].Machines[1])
			assert.Equal(t, tt.want[2].Machines[0], got[2].Machines[0])
		})
	}
}

func TestFindAvailableMachine(t *testing.T) {
	testRepo := FakeMachineRepo{}
	want := []string{"test-machine-1", "test-machine-2"}
	got, _ := FindAvailableMachine("test-gym-id", "test-machine-name", testRepo)
	assert.Equal(t, want, got)
}
