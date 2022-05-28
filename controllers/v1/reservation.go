package v1

import (
	"SADBackend/constant"
	"SADBackend/controllers/service"
	"SADBackend/model"
	"SADBackend/repo"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// period defintion:
// morning   06:00 ~ 12:00
// afternoon 12:00 ~ 18:00
// evening   18:00 ~ 21:00
// midnight  21:00 ~ 24:00
type timeRange string

const (
	TIME_MORNING   timeRange = "morning"
	TIME_AFTERNOON timeRange = "afternoon"
	TIME_EVENING   timeRange = "evening"
	TIME_MIDNIGHT  timeRange = "midnight"
)

type MakeReservationReq struct {
	Date      string `json:"date" example:"2006-01-02"`
	UserID    string `json:"account" example:"meowmeow123"`
	MachineID string `json:"machine_id" example:"machine-1000007"`
	Start     string `json:"start" example:"13:00"`
}

// @Summary Get Client Reservations
// @Produce json
// @Tags Client
// @Param account path string true "account e.g. meowmeow123"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/reservation/{account} [get]
func GetClientReservation(c *gin.Context, reservationDB repo.ReservationRepo) {
	userID := c.Param("account")

	var results []model.AggrReservationRes
	if err := reservationDB.GetReservation(userID, &results); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	res, _ := service.FindFutureKReservation(results, 4)

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

// @Summary Make Reservation
// @Produce json
// @Tags Client
// @Param MakeReservationReq  body MakeReservationReq true "date, account, machine_id, start"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/reservation [post]
func MakeReservation(c *gin.Context, reservationDB repo.ReservationRepo) {
	var mrReq MakeReservationReq
	if err := c.ShouldBindJSON(&mrReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	start, err := service.ComputeStartTime(mrReq.Date, mrReq.Start)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	err = reservationDB.Exist(mrReq.UserID, start, &struct{}{})
	if err == nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": "already has reservation at the same time"})
		return
	}

	if err := reservationDB.MakeReservation(mrReq.UserID, mrReq.MachineID, start); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}

// @Summary Get Available Time
// @Produce json
// @Tags Client
// @Param date query string true "format: 2006-01-02"
// @Param period query string true "period morning, afternoon, evening, midnight"
// @Param account query string true "e.g. meowmeow123"
// @Param branch_gym_id query string true "e.g. branch-1000001"
// @Param machine query string true "e.g. treadmill"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/available [get]
func GetAvailableTime(c *gin.Context, machineDB repo.MachineRepo, reservationDB repo.ReservationRepo) {
	date := c.Query("date")
	period := c.Query("period")
	userID := c.Query("account")
	gymID := c.Query("branch_gym_id")
	machinePrefix := c.Query("machine")

	// find existing reservations
	machineIDs, err := service.FindAvailableMachine(gymID, machinePrefix, machineDB)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	if len(machineIDs) == 0 {
		constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, machineIDs)
		return
	}

	lb, ub, err := getTimeRangeBound(period, date)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	var reservations []model.Reservation
	if err := reservationDB.QueryExistReservation(machineIDs, *lb, *ub, &reservations); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	res, err := service.FindKAvailabeTime(userID, reservations, machineIDs, *lb, *ub, 5)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

func getTimeRangeBound(tr string, date string) (*time.Time, *time.Time, error) {
	startBase, err := string2Time(date, "2006-01-02")
	if err != nil {
		return nil, nil, err
	}
	var ub, lb time.Time
	if tr == string(TIME_MORNING) {
		lb = startBase.Add(time.Duration(6) * time.Hour)
		ub = startBase.Add(time.Duration(12) * time.Hour).Add(-30 * time.Minute)
	} else if tr == string(TIME_AFTERNOON) {
		lb = startBase.Add(time.Duration(12) * time.Hour)
		ub = startBase.Add(time.Duration(18) * time.Hour).Add(-30 * time.Minute)
	} else if tr == string(TIME_EVENING) {
		lb = startBase.Add(time.Duration(18) * time.Hour)
		ub = startBase.Add(time.Duration(21) * time.Hour).Add(-30 * time.Minute)
	} else if tr == string(TIME_MIDNIGHT) {
		lb = startBase.Add(time.Duration(21) * time.Hour)
		ub = startBase.Add(time.Duration(24) * time.Hour).Add(-30 * time.Minute)
	} else {
		return nil, nil, fmt.Errorf("unknown time range: %s", tr)
	}
	return &lb, &ub, nil
}
