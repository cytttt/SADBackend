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

func TestMakeReservation(t *testing.T) {
	testRepo := repo.AllRepo{
		Reservation: FakeReservationRepo{},
	}
	gin.SetMode(gin.TestMode)
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqBody        v1.MakeReservationReq
		want           string
		wantHttpStatus int
	}{
		{
			name: "test case 1: happy",
			reqBody: v1.MakeReservationReq{
				Date:      "2022-02-04",
				Start:     "13:30",
				UserID:    testClientID,
				MachineID: testMachineID,
			},
			want:           "{\"code\":200,\"data\":null,\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bString, _ := json.Marshal(tt.reqBody)
			w, m := request(
				router,
				"POST",
				"/api/v1/user/reservation",
				bytes.NewBuffer(bString),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestGetClientReservation(t *testing.T) {
	testRepo := repo.AllRepo{
		Reservation: FakeReservationRepo{},
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
			name:           "test case 1: happy ",
			reqParam:       testClientID,
			want:           "{\"code\":200,\"data\":[{\"category\":\"abs\",\"date\":\"0001-01-01T00:00:00Z\",\"gym_id\":\"test-gym-1\",\"gym_name\":\"test-gym-daan\",\"machine_id\":\"test-machine-t1\",\"machine_name\":\"treadmill\"}],\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, m := request(
				router,
				"GET",
				"/api/v1/user/reservation/"+tt.reqParam,
				bytes.NewBuffer([]byte(``)),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestGetAvailableTime(t *testing.T) {
	testRepo := repo.AllRepo{
		Machine:     FakeMachineRepo{},
		Reservation: FakeReservationRepo{},
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
				"date":          "2023-09-09",
				"period":        "morning",
				"account":       testClientID,
				"branch_gym_id": testGymID,
				"machine":       "treadmill",
			},
			want:           "Ok",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, m := request(
				router,
				"GET",
				"/api/v1/user/available",
				bytes.NewBuffer([]byte(``)),
				tt.reqQuery,
			)
			// mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, m["msg"])
		})
	}
}
