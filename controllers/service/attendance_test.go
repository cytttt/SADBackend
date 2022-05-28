package service

import (
	"SADBackend/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostprocessStatData(t *testing.T) {
	tests := []struct {
		name    string
		args    []model.StatInSecond
		want    float32
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: integer",
			args: []model.StatInSecond{
				{
					AvgStaySecond: 3600,
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "test case 2: irrational number",
			args: []model.StatInSecond{
				{
					AvgStaySecond: 4800,
				},
			},
			want:    1.3333334,
			wantErr: false,
		},
		{
			name: "test case 3: rational number",
			args: []model.StatInSecond{
				{
					AvgStaySecond: 5400,
				},
			},
			want:    1.5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostprocessStatData(tt.args)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want, got[0].AvgStayTime)
		})
	}
}
