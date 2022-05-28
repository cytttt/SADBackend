package v1

import (
	"SADBackend/constant"
	"SADBackend/controllers/service"
	"SADBackend/model"
	"SADBackend/repo"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
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

// @Summary Client Signup
// @Produce json
// @Tags Client
// @Param signupCredentials body model.SignupReq true "account, password, name, email, gender, phone, birthday, height, weight"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/signup [post]
func Signup(c *gin.Context, clientDB repo.ClientRepo) {
	var signupReq model.SignupReq
	if err := c.ShouldBindJSON(&signupReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	if err := clientDB.Exist(signupReq.Account, struct{}{}); err == nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_EXISTS, gin.H{"error": err.Error()})
		return
	}

	newClient, err := service.PreprocessSignupInfo(signupReq)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}

	if err := clientDB.Signup(*newClient); err != nil {
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
func Login(c *gin.Context, clientDB repo.ClientRepo, staffDB repo.StaffRepo) {
	var loginCred LoginCred
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}
	if loginCred.UserRole == constant.USER_ROLE_Client {
		var client *model.Client
		if err := clientDB.Exist(loginCred.UserID, &client); err != nil {
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
		if err := staffDB.Exist(loginCred.UserID, &staff); err != nil {
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
func GetClientInfo(c *gin.Context, clientDB repo.ClientRepo) {
	userID := c.Query("account")

	var clientInfo *ClientInfoResp
	if err := clientDB.Exist(userID, &clientInfo); err != nil {
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
func GetClientStat(c *gin.Context, clientDB repo.ClientRepo) {
	userID := c.Param("account")
	var client *model.Client
	if err := clientDB.Exist(userID, &client); err != nil {
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
func UpdateClientInfo(c *gin.Context, clientDB repo.ClientRepo) {
	var updateReq model.UpdateUserInfoReq
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}

	err := clientDB.Exist(updateReq.Account, &struct{}{})
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, gin.H{"error": err.Error()})
		return
	}

	update := service.PreprocessUpdateInfo(updateReq)

	var clientInfo *ClientInfoResp
	if err := clientDB.UpdateClientInfo(updateReq.Account, update, &clientInfo); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, *clientInfo)
}

// @Summary Get Company Stat
// @Produce json
// @Tags Staff
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/staff/stat [get]
func GetCompanyStat(c *gin.Context, attendanceDB repo.AttendanceRepo) {
	var data []model.StatInSecond
	if err := attendanceDB.CompanyStat7days(&data); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	res, _ := service.PostprocessStatData(data)

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
