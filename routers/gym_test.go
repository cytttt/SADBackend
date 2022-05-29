package routes_test

import (
	"SADBackend/repo"
	routes "SADBackend/routers"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetGymList(t *testing.T) {
	testRepo := repo.AllRepo{
		Gym: FakeGymRepo{},
	}
	gin.SetMode(gin.TestMode)
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		want           string
		wantHttpStatus int
	}{
		{
			name:           "tset case 1: happy",
			want:           "{\"code\":200,\"data\":[{\"address\":\"\",\"available_machine\":{\"abs\":null,\"arm\":null,\"back\":null,\"cardio\":null,\"chest\":null,\"hips\":null,\"leg\":null},\"branch_gym_id\":\"test-gym-1\",\"name\":\"test-gym-daan\",\"status\":\"uncrowded\"}],\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, m := request(
				router,
				"GET",
				"/api/v1/gym/list",
				bytes.NewBuffer([]byte(``)),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}
