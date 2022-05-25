package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type GetAvailableTimeResp struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	MachineID string `json:"machine_id"`
}

// @Summary Get Client Reservations
// @Produce json
// @Tags Client
// @Param account path string true "account e.g. meowmeow123"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/reservation/{account} [get]
func GetClientReservation(c *gin.Context) {
	account := c.Param("account")
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))

	matchStage := bson.M{
		"$match": bson.M{
			"user_id":  bson.M{"$eq": account},
			"start_at": bson.M{"$gte": time.Now().In(loc)},
		},
	}
	lookupStage1 := bson.M{
		"$lookup": bson.M{
			"from":         "machine",
			"localField":   "machine_id",
			"foreignField": "machine_id",
			"as":           "machines",
		},
	}
	lookupStage2 := bson.M{
		"$lookup": bson.M{
			"from":         "gym",
			"localField":   "machines.0.gym_id",
			"foreignField": "branch_gym_id",
			"as":           "gyms",
		},
	}
	pip := []bson.M{matchStage, lookupStage1, lookupStage2}
	cursor, err := mongodb.ReservationCollection.Aggregate(context.Background(), pip)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	var results []struct {
		ID          primitive.ObjectID `bson:"_id"`
		UserID      string             `bson:"user_id"`
		MachineID   string             `bson:"machine_id"`
		Category    model.PartCategory `bson:"category"`
		MachineName string             `bson:"machine_name"`
		StartAt     time.Time          `bson:"start_at"`
		Expired     bool               `bson:"expired"`
		Gyms        []model.BranchGym  `bson:"gyms"`
		Machines    []model.Machine    `bson:"machines"`
	}
	if err := cursor.All(context.TODO(), &results); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	var res []ReservationResp
	for _, i := range results {
		res = append(res, ReservationResp{
			MachineID:   i.MachineID,
			Category:    string(i.Machines[0].Category),
			MachineName: i.Machines[0].Name,
			GymID:       i.Gyms[0].BranchGymID,
			GymName:     i.Gyms[0].Name,
			Date:        i.StartAt,
		})
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

// @Summary Make Reservation
// @Produce json
// @Tags Client
// @Param MakeReservationReq  body MakeReservationReq true "date, account, machine_id, start"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/reservation [post]
func MakeReservation(c *gin.Context) {
	var mrReq MakeReservationReq
	if err := c.ShouldBindJSON(&mrReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}
	startBase, err := string2Time(mrReq.Date, "2006-01-02")
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	startHour, err := time.Parse("15:04", mrReq.Start)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	start := startBase.Add(time.Duration(startHour.Hour()) * time.Hour)
	newReservation := model.Reservation{
		ID:        primitive.NewObjectID(),
		UserID:    mrReq.UserID,
		MachineID: mrReq.MachineID,
		StartAt:   start,
	}
	if _, err := mongodb.ReservationCollection.InsertOne(context.Background(), newReservation); err != nil {
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
func GetAvailableTime(c *gin.Context) {
	date := c.Query("date")
	period := c.Query("period")
	userID := c.Query("account")
	gymID := c.Query("branch_gym_id")
	machinePrefix := c.Query("machine")
	// find machines
	machineFilter := bson.M{
		"reservation_only": true,
		"gym_id":           gymID,
		"name": primitive.Regex{
			Pattern: "^" + machinePrefix,
			Options: "",
		},
	}

	cursor, err := mongodb.MachineCollection.Find(context.Background(), machineFilter)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	var machines []model.Machine
	if err := cursor.All(context.TODO(), &machines); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	// find existing reservations
	var machineIDs []string
	for _, i := range machines {
		machineIDs = append(machineIDs, i.MachineID)
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
	reservationFilter := bson.M{
		"machine_id": bson.M{
			"$in": machineIDs,
		},
		"start_at": bson.M{"$gte": lb, "$lt": ub},
	}
	cursor, err = mongodb.ReservationCollection.Find(context.Background(), reservationFilter)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	var reservations []model.Reservation
	if err := cursor.All(context.TODO(), &reservations); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	results, err := findResult_K(userID, reservations, machineIDs, *lb, *ub, 5)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, results)
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

func findResult_K(userID string, res []model.Reservation, mID []string, lbd time.Time, ubd time.Time, k int) ([]GetAvailableTimeResp, error) {
	var tmpAvailableTime []GetAvailableTimeResp
	// setup map
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	reservationMap := make(map[string][]model.Reservation)
	for _, r := range res {
		tmpStr := r.StartAt.In(loc).Format("15:04")
		reservationMap[tmpStr] = append(reservationMap[tmpStr], r)
	}

	for d := lbd; d.After(ubd) == false; d = d.Add(30 * time.Minute) {
		st := d.Format("15:04")
		ed := d.Add(30 * time.Minute).Format("15:04")
		var mIDAtd []string // machine available at time d
		// if the user has already  made reservation at the same time => continue
		if _, ok := reservationMap[st]; ok {
			reserveConflict := false
			tmpSet := mapset.NewSet()
			for _, r := range reservationMap[st] {
				if r.UserID == userID {
					reserveConflict = true
					break
				}
				tmpSet.Add(r.MachineID)
			}
			if reserveConflict {
				continue
			}
			for _, m := range mID {
				if tmpSet.Contains(m) {
					continue
				}
				mIDAtd = append(mIDAtd, m)
			}
		} else {
			mIDAtd = append(mIDAtd, mID...)
		}
		tmpAvailableTime = append(tmpAvailableTime, GetAvailableTimeResp{
			Start:     st,
			End:       ed,
			MachineID: mIDAtd[rand.Intn(len(mIDAtd))],
		})
	}
	permutation := rand.Perm(len(tmpAvailableTime))[:Min(k, len(tmpAvailableTime))]
	sort.Ints(permutation)
	var AvailableTime []GetAvailableTimeResp
	for _, p := range permutation {
		AvailableTime = append(AvailableTime, tmpAvailableTime[p])
	}
	return AvailableTime, nil
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
