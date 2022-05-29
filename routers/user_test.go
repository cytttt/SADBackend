package routes_test

import (
	"SADBackend/constant"
	v1 "SADBackend/controllers/v1"
	"SADBackend/repo"
	routes "SADBackend/routers"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientInfo(t *testing.T) {
	testRepo := repo.AllRepo{
		Client: FakeClientRepo{},
	}
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqQuery       map[string]string
		want           string
		wantHttpStatus int
	}{
		{
			name: "test case 1:  happy",
			reqQuery: map[string]string{
				"account": testClientID,
			},
			want:           "{\"code\":200,\"data\":{\"account\":\"test-client-1\",\"body_info\":{\"Height\":0,\"Weight\":0},\"created_at\":\"0001-01-01T00:00:00Z\",\"email\":\"\",\"name\":\"test-client-amy\",\"payment_method\":{\"Account\":\"\",\"PayType\":\"\"},\"personal_info\":{\"Birthday\":\"0001-01-01T00:00:00Z\",\"Gender\":\"female\",\"Phone\":\"\"},\"subscription\":{\"ExpiredAt\":\"0001-01-01T00:00:00Z\",\"Plan\":\"\"},\"updated_at\":\"0001-01-01T00:00:00Z\"},\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w, m := request(
				router,
				"GET",
				"/api/v1/user/info",
				bytes.NewBuffer([]byte(``)),
				tt.reqQuery,
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestGetClientStat(t *testing.T) {
	testRepo := repo.AllRepo{
		Client: FakeClientRepo{},
	}
	router := routes.InitRouters(testRepo)
	tests := []struct {
		name           string
		reqParam       string
		want           string
		wantHttpStatus int
	}{
		{
			name:           "test case 1:  happy",
			reqParam:       testClientID,
			want:           "{\"code\":200,\"data\":{\"calories\":100,\"least_train\":\"\",\"most_train\":\"cardio\",\"stay_time\":0},\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w, m := request(
				router,
				"GET",
				"/api/v1/user/stat/"+tt.reqParam,
				bytes.NewBuffer([]byte(``)),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}

func TestLogin(t *testing.T) {
	testRepo := repo.AllRepo{
		Client: FakeClientRepo{},
		Staff:  FakeStaffRepo{},
	}
	router := routes.InitRouters(testRepo)
	// test client
	tests := []struct {
		name           string
		reqBody        v1.LoginCred
		want           string
		wantHttpStatus int
	}{
		{
			name: "test case 1: happy client",
			reqBody: v1.LoginCred{
				UserID:   testClientID,
				Password: testPwd,
				UserRole: constant.USER_ROLE_Client,
			},
			want:           "{\"code\":200,\"data\":{\"account\":\"test-client-1\",\"level\":\"\",\"name\":\"test-client-amy\",\"user_role\":\"client\"},\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
		{
			name: "test case 2: happy staff",
			reqBody: v1.LoginCred{
				UserID:   testStaffID,
				Password: testPwd,
				UserRole: constant.USER_ROLE_Staff,
			},
			want:           "{\"code\":200,\"data\":{\"account\":\"test-staff-1\",\"level\":\"rookie\",\"name\":\"test-staff-amy\",\"user_role\":\"staff\"},\"msg\":\"Ok\"}",
			wantHttpStatus: 200,
		},
		{
			name: "test case 3: client incorrect pwd",
			reqBody: v1.LoginCred{
				UserID:   testClientID,
				Password: "123",
				UserRole: constant.USER_ROLE_Client,
			},
			want:           "{\"code\":10004,\"data\":null,\"msg\":\"Incorrect password\"}",
			wantHttpStatus: 200,
		},
		{
			name: "test case 4: bad request missing user role",
			reqBody: v1.LoginCred{
				UserID:   testClientID,
				Password: testPwd,
			},
			want:           "{\"code\":400,\"data\":{\"error\":\"Key: 'LoginCred.UserRole' Error:Field validation for 'UserRole' failed on the 'required' tag\"},\"msg\":\"Invalid params error\"}",
			wantHttpStatus: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bString, _ := json.Marshal(tt.reqBody)
			w, m := request(
				router,
				"POST",
				"/api/v1/user/login",
				bytes.NewBuffer(bString),
				map[string]string{},
			)
			mJson, _ := json.Marshal(m)
			assert.Equal(t, tt.wantHttpStatus, w.Code)
			assert.Equal(t, tt.want, string(mJson))
		})
	}
}
