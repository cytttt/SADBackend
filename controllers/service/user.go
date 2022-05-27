package service

import (
	"SADBackend/constant"
	"SADBackend/model"
	"crypto/sha256"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PreprocessSignupInfo(req model.SignupReq) (*model.Client, error) {
	pwdCrypto := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))
	birthdayTime, err := string2Time(req.Birthday, "2006/01/02")
	if err != nil {
		return nil, err
	}
	newClient := model.Client{
		ID:       primitive.NewObjectID(),
		UserID:   req.Account,
		Name:     req.Name,
		Email:    req.Email,
		Password: pwdCrypto,
		PersonalInfo: model.UserInfo{
			Gender:   req.Gender,
			Phone:    req.Phone,
			Birthday: *birthdayTime,
		},
		BodyInfo: model.BodyInfo{
			Weight: float64(req.Weight),
			Height: float64(req.Height),
		},
	}
	return &newClient, nil
}

func PreprocessUpdateInfo(req model.UpdateUserInfoReq) bson.M {
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	update := bson.M{
		"$set": bson.M{
			"name":                    req.Name,
			"email":                   req.Email,
			"personal_info.gender":    req.Gender,
			"personal_info.phone":     req.Phone,
			"personal_info.birthday":  time.Date(req.Year, time.Month(req.Month), req.Day, 0, 0, 0, 0, loc),
			"body_info.weight":        req.Weight,
			"body_info.height":        req.Height,
			"subscription.plan":       req.Plan,
			"payment_method.pay_type": req.PayType,
			"payment_method.account":  req.PaymentAccount,
			"updated_at":              time.Now().In(loc),
		},
	}
	return update
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

func VerifyPwd(req, ref string) error {
	pwdCrypto := fmt.Sprintf("%x", sha256.Sum256([]byte(req)))
	if pwdCrypto != ref {
		return constant.NewError(constant.ERROR_INCORRECT_PASSWORD)
	}
	return nil
}
