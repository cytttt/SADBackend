package service

import (
	"SADBackend/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostprocessGymList(t *testing.T) {
	tests := []struct {
		name    string
		args    []model.BranchGym
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: crowded",
			args: []model.BranchGym{
				{
					CurrentNumberPeople: 90,
					Info: model.BranchInfo{
						ClientNumberLimit: 100,
					},
				},
			},
			want:    "crowded",
			wantErr: false,
		},
		{
			name: "test case 2: crowded edge case",
			args: []model.BranchGym{
				{
					CurrentNumberPeople: 80,
					Info: model.BranchInfo{
						ClientNumberLimit: 100,
					},
				},
			},
			want:    "uncrowded",
			wantErr: false,
		},
		{
			name: "test case 3: uncrowded",
			args: []model.BranchGym{
				{
					CurrentNumberPeople: 70,
					Info: model.BranchInfo{
						ClientNumberLimit: 100,
					},
				},
			},
			want:    "uncrowded",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostprocessGymList(tt.args)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want, got[0].Status)
		})
	}
}
