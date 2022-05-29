package routes_test

import (
	v1 "SADBackend/controllers/v1"
	"SADBackend/model"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	testPwd          = "test-pwd"
	testClientID     = "test-client-1"
	testClientName   = "test-client-amy"
	testClientGender = "female"
	testStaffID      = "test-staff-1"
	testStaffName    = "test-staff-amy"
	testStaffLevel   = "rookie"
)

type FakeClientRepo struct{}

func (_ FakeClientRepo) Exist(userID string, result interface{}) error {
	switch result.(type) {
	case **model.Client:
		testPwdHashed := fmt.Sprintf("%x", sha256.Sum256([]byte(testPwd)))
		res := result.(**model.Client)
		*res = &model.Client{
			Password: testPwdHashed,
			UserID:   testClientID,
			Name:     testClientName,
			Statistics: model.Stat{
				Calories:  100,
				MostTrain: model.PART_CARDIO,
			},
		}
	case **v1.ClientInfoResp:
		res := result.(**v1.ClientInfoResp)
		*res = &v1.ClientInfoResp{
			UserID: testClientID,
			Name:   testClientName,
			PersonalInfo: model.UserInfo{
				Gender: testClientGender,
			},
		}
	default:
		log.Println("check exist")
	}
	return nil
}

func (_ FakeClientRepo) UpdateClientInfo(userID string, update bson.M, result interface{}) error {
	return nil
}
func (_ FakeClientRepo) Signup(newClient model.Client) error { return nil }

type FakeStaffRepo struct{}

func (_ FakeStaffRepo) Exist(userID string, result interface{}) error {
	switch result.(type) {
	case **model.Staff:
		testPwdHashed := fmt.Sprintf("%x", sha256.Sum256([]byte(testPwd)))
		res := result.(**model.Staff)
		*res = &model.Staff{
			Password: testPwdHashed,
			UserID:   testStaffID,
			Name:     testStaffName,
			Level:    testStaffLevel,
		}
	default:
		log.Println("unknown result type")
	}
	return nil
}

type FakeReservationRepo struct{}
type FakeGymRepo struct{}
type FakeMachineRepo struct{}
type FakeAttendanceRepo struct{}

func request(router *gin.Engine, method string, path string, body *bytes.Buffer, query map[string]string) (*httptest.ResponseRecorder, map[string]interface{}) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	for _, i := range query {
		q.Add(i, query[i])
	}
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	m := make(map[string]interface{})
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	return w, m
}
