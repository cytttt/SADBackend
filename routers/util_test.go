package routes_test

import (
	"SADBackend/constant"
	v1 "SADBackend/controllers/v1"
	"SADBackend/model"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	testPwd             = "test-pwd"
	testClientID        = "test-client-1"
	testClientID2       = "test-client-2"
	testClientName      = "test-client-amy"
	testClientEmail     = "client@test.com"
	testClientGender    = "female"
	testStaffID         = "test-staff-1"
	testStaffName       = "test-staff-amy"
	testStaffLevel      = "rookie"
	testGymID           = "test-gym-1"
	testGymName         = "test-gym-daan"
	testMachineID       = "test-machine-t1"
	testMachineName     = "test-machine-treadmill"
	testMachineCategory = model.PART_CARDIO
)

type FakeClientRepo struct{}
type FakeStaffRepo struct{}
type FakeReservationRepo struct{}
type FakeGymRepo struct{}
type FakeMachineRepo struct{}
type FakeAttendanceRepo struct{}

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
		if userID != testClientID {
			return constant.NewError(constant.ERROR)
		}
		log.Println("check exist")
	}
	return nil
}

func (_ FakeClientRepo) UpdateClientInfo(userID string, update bson.M, result interface{}) error {
	switch result.(type) {
	case **v1.ClientInfoResp:
		res := result.(**v1.ClientInfoResp)
		*res = &v1.ClientInfoResp{
			UserID: testClientID,
		}
	default:
		if userID != testClientID {
			return constant.NewError(constant.ERROR)
		}
		log.Println("check exist")
	}
	return nil
}
func (_ FakeClientRepo) Signup(newClient model.Client) error { return nil }

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

func (_ FakeReservationRepo) Exist(userID string, startAt time.Time, result interface{}) error {
	return constant.NewError(constant.ERROR)
}

func (_ FakeReservationRepo) GetReservation(userID string, results interface{}) error {
	switch results.(type) {
	case *[]model.AggrReservationRes:
		res := results.(*[]model.AggrReservationRes)
		*res = append(*res, model.AggrReservationRes{
			UserID:    testClientID,
			MachineID: testMachineID,
			Machines: []model.Machine{
				{
					Name:     "treadmill",
					Category: model.PART_ABS,
				},
			},
			Gyms: []model.BranchGym{
				{
					BranchGymID: testGymID,
					Name:        testGymName,
				},
			},
		})
	}
	return nil
}

func (_ FakeReservationRepo) MakeReservation(userID, machineID string, start time.Time) error {
	return nil
}

func (_ FakeReservationRepo) QueryExistReservation(machineIDs []string, lb, ub time.Time, results interface{}) error {
	return nil
}

func (_ FakeGymRepo) GymList(results interface{}) error {
	switch results.(type) {
	case *[]model.BranchGym:
		res := results.(*[]model.BranchGym)
		*res = append(*res, model.BranchGym{
			BranchGymID:         testGymID,
			Name:                testGymName,
			CurrentNumberPeople: 80,
			Info: model.BranchInfo{
				ClientNumberLimit: 100,
			},
		})
	default:
		log.Print("unknown type")
	}
	return nil
}

func (_ FakeMachineRepo) MachineList(gymID string, results interface{}) error {
	switch results.(type) {
	case *[]model.Machine:
		res := results.(*[]model.Machine)
		*res = append(*res, model.Machine{
			MachineID:  testMachineID,
			Name:       testMachineName,
			Category:   testMachineCategory,
			WaitingPPL: 10,
		})
	default:
		log.Print("unknown type")
	}
	return nil
}

func (_ FakeMachineRepo) UpdateAmount(machineID string, amount int, result interface{}) error {
	switch result.(type) {
	case **model.Machine:
		res := result.(**model.Machine)
		*res = &model.Machine{
			MachineID:  testMachineID,
			Name:       testMachineName,
			Category:   testMachineCategory,
			WaitingPPL: 10,
		}
	default:
		log.Println("unknown result type")
	}
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

func (_ FakeAttendanceRepo) CompanyStat7days(results interface{}) error {
	switch results.(type) {
	case *[]model.StatInSecond:
		res := results.(*[]model.StatInSecond)
		*res = append(*res, model.StatInSecond{
			AttendanceCount: 10,
			AvgStaySecond:   3600,
		})
		*res = append(*res, model.StatInSecond{
			AttendanceCount: 12,
			AvgStaySecond:   4800,
		})
	default:
		log.Print("unknown type")
	}
	return nil
}

func request(router *gin.Engine, method string, path string, body *bytes.Buffer, query map[string]string) (*httptest.ResponseRecorder, map[string]interface{}) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	for k, i := range query {
		q.Add(k, i)
	}
	req.URL.RawQuery = q.Encode()
	log.Println("@@", req.URL.String())
	router.ServeHTTP(w, req)

	m := make(map[string]interface{})
	_ = json.Unmarshal(w.Body.Bytes(), &m)
	return w, m
}
