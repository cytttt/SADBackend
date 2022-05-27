package service

import (
	"SADBackend/constant"
	"SADBackend/model"
	"crypto/sha256"
	"fmt"
	"time"

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
