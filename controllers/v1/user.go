package v1

import (
	"SADBackend/constant"
	"SADBackend/controllers/service"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"SADBackend/repo"
	"context"
	"net/http"
	"net/mail"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginCred struct {
	UserID   string            `json:"account" binding:"required" example:"meowmeow123"`
	Password string            `json:"password" binding:"required" example:"meowmoew22"`
	UserRole constant.UserRole `json:"user_role" binding:"required" example:"client"`
}

type LoginResp struct {
	Account  string            `json:"account"`
	Name     string            `json:"name"`
	UserRole constant.UserRole `json:"user_role"`
	Level    string            `json:"level"`
}

type ClientInfoResp struct {
	UserID       string                 `bson:"user_id" json:"account"`
	Name         string                 `bson:"name" json:"name"`
	Email        string                 `bson:"email" json:"email"`
	PersonalInfo model.UserInfo         `bson:"personal_info" json:"personal_info"`
	BodyInfo     model.BodyInfo         `bson:"body_info" json:"body_info"`
	Subscription model.SubscriptionInfo `bson:"subscription" json:"subscription"`
	Payment      model.PaymentMethod    `bson:"payment_method" json:"payment_method"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at" `
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}

type ReservationResp struct {
	Category    string    `json:"category"`
	MachineID   string    `json:"machine_id"`
	MachineName string    `json:"machine_name"`
	GymID       string    `json:"gym_id"`
	GymName     string    `json:"gym_name"`
	Date        time.Time `json:"date"`
}

type CompanyStatResp struct {
	Date        string  `json:"date"`
	Attendance  int     `json:"attendance_count"`
	AvgStayTime float32 `json:"avg_stay_hour"`
}

// @Summary Client Signup
// @Produce json
// @Tags Client
// @Param signupCredentials body model.SignupReq true "account, password, name, email, gender, phone, birthday, height, weight"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/signup [post]
func Signup(c *gin.Context) {
	var signupReq model.SignupReq
	if err := c.ShouldBindJSON(&signupReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	if err := repo.Client.Exist(signupReq.Account, struct{}{}); err == nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_EXISTS, gin.H{"error": err.Error()})
		return
	}

	newClient, err := service.PreprocessSignupInfo(signupReq)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	if err := repo.Client.Signup(*newClient); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, nil)
}

// @Summary User Login
// @Produce json
// @Tags User
// @Param loginCredentials  body LoginCred true "account only, password, userRole("client","staff")"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	var loginCred LoginCred
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}
	if loginCred.UserRole == constant.USER_ROLE_Client {
		var client *model.Client
		if err := repo.Client.Exist(loginCred.UserID, &client); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, gin.H{"error": err.Error()})
			return
		}
		if err := service.VerifyPwd(loginCred.Password, client.Password); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_INCORRECT_PASSWORD, nil)
			return
		}
		constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, LoginResp{
			Account:  client.UserID,
			Name:     client.Name,
			UserRole: constant.USER_ROLE_Client,
		})
		return
	} else if loginCred.UserRole == constant.USER_ROLE_Staff {
		var staff *model.Staff
		if err := repo.Staff.Exist(loginCred.UserID, &staff); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, gin.H{"error": err.Error()})
			return
		}
		if err := service.VerifyPwd(loginCred.Password, staff.Password); err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_INCORRECT_PASSWORD, nil)
			return
		}
		constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, LoginResp{
			Account:  staff.UserID,
			Name:     staff.Name,
			UserRole: constant.USER_ROLE_Staff,
			Level:    staff.Level,
		})
		return
	}
	constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": "invalid user role"})
}

// @Summary Get Client Info
// @Produce json
// @Tags Client
// @param account query string true "account e.g. meowmeow123"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/info [get]
func GetClientInfo(c *gin.Context) {
	userID := c.Query("account")

	var clientInfo *ClientInfoResp
	if err := repo.Client.Exist(userID, &clientInfo); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_EXISTS, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, *clientInfo)
}

// @Summary Get Client Info
// @Produce json
// @Tags Client
// @Param account path string true "account e.g. meowmeow123"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/stat/{account} [get]
func GetClientStat(c *gin.Context) {
	userID := c.Param("account")
	var client *model.Client
	if err := repo.Client.Exist(userID, &client); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_EXISTS, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, client.Statistics)
}

// @Summary Update Client info
// @Produce json
// @Tags Client
// @Param UpdateClientInfo body model.UpdateUserInfoReq true "account, ..."
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/info [put]
func UpdateClientInfo(c *gin.Context) {
	var updateReq model.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}
	// check use exist or not
	err := mongodb.ClientCollection.FindOne(context.Background(), bson.M{"user_id": updateReq.Account}).Decode(&struct{}{})
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	// update
	//opt := &options.UpdateOptions{}
	//opt.SetUpsert(true)
	opt := options.FindOneAndUpdate()
	opt.SetUpsert(true)
	opt.SetReturnDocument(options.After)
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	filter := bson.M{"user_id": updateReq.Account}
	update := bson.M{
		"$set": bson.M{
			"name":                    updateReq.Name,
			"email":                   updateReq.Email,
			"personal_info.gender":    updateReq.Gender,
			"personal_info.phone":     updateReq.Phone,
			"personal_info.birthday":  time.Date(updateReq.Year, time.Month(updateReq.Month), updateReq.Day, 0, 0, 0, 0, loc),
			"body_info.weight":        updateReq.Weight,
			"body_info.height":        updateReq.Height,
			"subscription.plan":       updateReq.Plan,
			"payment_method.pay_type": updateReq.PayType,
			"payment_method.account":  updateReq.PaymentAccount,
			"updated_at":              time.Now().In(loc),
		},
	}
	var clientInfo ClientInfoResp
	if err := mongodb.ClientCollection.FindOneAndUpdate(context.Background(), filter, update, opt).Decode(&clientInfo); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, clientInfo)
}

// @Summary Get Company Stat
// @Produce json
// @Tags Staff
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/staff/stat [get]
func GetCompanyStat(c *gin.Context) {
	// 理論上 loc should be passed
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	cur := time.Now().In(loc)
	y, m, d := cur.Date()
	ub := time.Date(y, m, d, 0, 0, 0, 0, loc)
	lb := ub.AddDate(0, 0, -7)
	matchStage := bson.M{
		"$match": bson.M{
			"enter": bson.M{"$gte": lb, "$lt": ub},
		},
	}
	groupStage := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y/%m/%d",
					"date":   "$enter",
				}},
			"attendance_count": bson.M{"$sum": 1},
			"avg_stay_second":  bson.M{"$avg": "$stay_time"},
		},
	}
	pip := []bson.M{matchStage, groupStage}
	cursor, err := mongodb.AttendanceCollection.Aggregate(context.Background(), pip)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	var results []struct {
		Date            string  `bson:"_id"`
		AttendanceCount int     `bson:"attendance_count"`
		AvgStaySecond   float32 `bson:"avg_stay_second"`
	}
	if err := cursor.All(context.TODO(), &results); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date < results[j].Date
	})
	var res []CompanyStatResp
	for _, i := range results {
		res = append(res, CompanyStatResp{
			Date:        i.Date,
			Attendance:  i.AttendanceCount,
			AvgStayTime: i.AvgStaySecond / 3600,
		})
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, res)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func string2Time(timeStr, format string) (*time.Time, error) {
	offset := int((8 * time.Hour).Seconds())
	loc := time.FixedZone("Asia/Taipei", offset)
	newTime, err := time.ParseInLocation(format, timeStr, loc)
	if err != nil {
		return nil, err
	}
	return &newTime, err
}

func deleteDocument(account string) error {
	if _, err := mongodb.ClientCollection.DeleteOne(context.Background(), bson.M{"user_id": account}); err != nil {
		return err
	}
	return nil
}
