package routes_test

import (
	v1 "SADBackend/controllers/v1"
	"SADBackend/repo"
	routes "SADBackend/routers"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMachineList(t *testing.T) {
	testRepo := repo.AllRepo{
		Machine: FakeMachineRepo{},
	}
	gin.SetMode(gin.TestMode)
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqQuery       map[string]string
		want           string
		wantHttpStatus int
	}{
		{
			name: "test case 1: happy",
			reqQuery: map[string]string{
				"gym_id": testGymID,
			},
			wantHttpStatus: 200,
			want:           "{\"code\":200,\"data\":[{\"category\":\"cardio\",\"machine_id\":\"test-machine-t1\",\"name\":\"test-machine-treadmill\",\"waiting_ppl\":10}],\"msg\":\"Ok\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, m := request(
				router,
				"GET",
				"/api/v1/gym/machine",
				bytes.NewBuffer([]byte(``)),
				tt.reqQuery,
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestGetMachineListByCategory(t *testing.T) {
	testRepo := repo.AllRepo{
		Machine: FakeMachineRepo{},
	}
	gin.SetMode(gin.TestMode)
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqParam       string
		want           string
		wantHttpStatus int
	}{
		{
			name:           "test case 1: happy",
			reqParam:       testGymID,
			wantHttpStatus: 200,
			want:           "{\"code\":200,\"data\":[{\"category\":\"back\",\"machines\":null},{\"category\":\"chest\",\"machines\":null},{\"category\":\"cardio\",\"machines\":[{\"category\":\"cardio\",\"machine_id\":\"test-machine-t1\",\"name\":\"test-machine-treadmill\",\"waiting_ppl\":10}]},{\"category\":\"abs\",\"machines\":null},{\"category\":\"leg\",\"machines\":null},{\"category\":\"arm\",\"machines\":null},{\"category\":\"hips\",\"machines\":null}],\"msg\":\"Ok\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, m := request(
				router,
				"GET",
				"/api/v1/gym/machine/category/"+tt.reqParam,
				bytes.NewBuffer([]byte(``)),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestUpdateMachineStatus(t *testing.T) {
	testRepo := repo.AllRepo{
		Machine: FakeMachineRepo{},
	}
	gin.SetMode(gin.TestMode)
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqBody        v1.UpdateMachineStatusReq
		want           string
		wantHttpStatus int
	}{
		{
			name: "test case 1: happy",
			reqBody: v1.UpdateMachineStatusReq{
				MachineID: testMachineID,
				Amount:    1,
			},
			wantHttpStatus: 200,
			want:           "{\"code\":200,\"data\":[{\"category\":\"cardio\",\"machine_id\":\"test-machine-t1\",\"name\":\"test-machine-treadmill\",\"waiting_ppl\":10}],\"msg\":\"Ok\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bString, _ := json.Marshal(tt.reqBody)
			w, m := request(
				router,
				"PUT",
				"/api/v1/machine/status",
				bytes.NewBuffer(bString),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}
