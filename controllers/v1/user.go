package v1

import (
	"SADBackend/constant"
	"SADBackend/model"
	"SADBackend/pkg/mongodb"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginCred struct {
	Account  string            `json:"account" binding:"required" example:"meowmeow123"`
	Password string            `json:"password" binding:"required" example:"meowmoew22"`
	UserRole constant.UserRole `json:"user_role" binding:"required" example:"client"`
}

type LoginResp struct {
	Account  string            `json:"account"`
	Name     string            `json:"name"`
	UserRole constant.UserRole `json:"user_role"`
	Level    string            `json:"level"`
}

type UpdateUserInfoReq struct {
	Account string  `json:"account" example:"meowmeow123"` // use to identify user
	Name    string  `json:"name" example:"testMeowClient"`
	Email   string  `json:"email" binding:"email" example:"meowtestclient@gmail.com"`
	Gender  string  `json:"gender" example:"male"`
	Phone   string  `json:"phone" example:"0919886886"`
	Year    int     `json:"year" example:"2001"`
	Month   int     `json:"month" example:"5"`
	Day     int     `json:"day" example:"29"`
	Weight  float64 `json:"weihgt" example:"69.69"`
	Height  float64 `json:"height" example:"180.13"`
}
type ClientInfoResp struct {
	UserID           string                 `bson:"user_id" json:"account"`
	Name             string                 `bson:"name" json:"name"`
	Email            string                 `bson:"email" json:"email"`
	PersonalInfo     model.UserInfo         `bson:"personal_info" json:"personal_info"`
	BodyInfo         model.BodyInfo         `bson:"body_info" json:"body_info"`
	Subscription     model.SubscriptionInfo `bson:"subscription" json:"subscription"`
	Payment          model.PaymentMethod    `bson:"payment_method" json:"payment_method"`
	AttendenceRecord []model.Attendence     `bosn:"attendence_record" josn:"attendence_record"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at" `
	UpdatedAt        time.Time              `bson:"updated_at" json:"updated_at"`
}

// @Summary User Login
// @Produce json
// @Tags User
// @Param loginCredentials  body LoginCred true "account/email, password, userRole("client","staff")"
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	var loginCred LoginCred
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, gin.H{"error": err.Error()})
		return
	}
	filter := bson.M{}
	if validEmail(loginCred.Account) {
		filter = bson.M{"email": loginCred.Account}
	} else {
		filter = bson.M{"user_id": loginCred.Account}
	}

	if loginCred.UserRole == constant.USER_ROLE_Client {
		var client model.Client
		passwdCrypto := fmt.Sprintf("%x", sha256.Sum256([]byte(loginCred.Password)))
		err := mongodb.ClientCollection.FindOne(context.Background(), filter).Decode(&client)
		if err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, nil)
			return
		}
		if passwdCrypto != client.Password {
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
		var staff model.Staff
		passwdCrypto := fmt.Sprintf("%x", sha256.Sum256([]byte(loginCred.Password)))
		err := mongodb.StaffCollection.FindOne(context.Background(), filter).Decode(&staff)
		if err != nil {
			constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, nil)
			return
		}
		if passwdCrypto != staff.Password {
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
	constant.ResponseWithData(c, http.StatusBadRequest, constant.INVALID_PARAMS, "invalid user role")
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
	var clientInfo ClientInfoResp
	err := mongodb.ClientCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&clientInfo)
	if err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR_USER_NOT_FOUND, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, clientInfo)
}

// @Summary Update Client info
// @Produce json
// @Tags Client
// @Param UpdateClientInfo body UpdateUserInfoReq true "account, ..."
// @Success 200 {object} constant.Response
// @Failure 500 {object} constant.Response
// @Router /api/v1/user/info [put]
func UpdateClientInfo(c *gin.Context) {
	var updateReq UpdateUserInfoReq
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
			"name":                   updateReq.Name,
			"email":                  updateReq.Email,
			"personal_info.gender":   updateReq.Gender,
			"personal_info.phone":    updateReq.Phone,
			"personal_info.birthday": time.Date(updateReq.Year, time.Month(updateReq.Month), updateReq.Day, 0, 0, 0, 0, loc),
			"body_info.weight":       updateReq.Weight,
			"body_info.height":       updateReq.Height,
			"updated_at":             time.Now().In(loc),
		},
	}
	var clientInfo ClientInfoResp
	if err := mongodb.ClientCollection.FindOneAndUpdate(context.Background(), filter, update, opt).Decode(&clientInfo); err != nil {
		constant.ResponseWithData(c, http.StatusOK, constant.ERROR, gin.H{"error": err.Error()})
		return
	}
	constant.ResponseWithData(c, http.StatusOK, constant.SUCCESS, clientInfo)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
