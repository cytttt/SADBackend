package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	PersonalInfo UserInfo           `bson:"personal_info"`
	BodyInfo     BodyInfo           `bson:"body_info"`
	Subscription SubscriptionInfo   `bson:"subscription"`
	Payment      PaymentMethod      `bson:"payment_method"`
	Statistics   Stat               `bson:"statistics"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type Staff struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	PersonalInfo UserInfo           `bson:"personal_info"`
	Level        string             `bson:"level"`
}

type UserInfo struct {
	Gender   string    `bson:"gender"`
	Phone    string    `bson:"phone"`
	Birthday time.Time `bson:"birthday"`
}

type BodyInfo struct {
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type SubscriptionInfo struct {
	Plan      PlanLevel `bson:"plan"`
	ExpiredAt time.Time `bson:"expired_at"`
}

type PaymentMethod struct {
	PayType string `bson:"pay_type"`
	Account string `bson:"account"`
}

type Stat struct {
	StayTime   float64      `bson:"stay_time" json:"stay_time"`
	Calories   int          `bson:"calories" json:"calories"`
	MostTrain  PartCategory `bson:"most_train" json:"most_train"`
	LeastTrain PartCategory `bson:"least_train" json:"least_train"`
}

type PlanLevel string

const (
	PLAN_Basic    PlanLevel = "Basic"
	PLAN_Standard PlanLevel = "Standard"
	PLAN_Premium  PlanLevel = "Premium"
)

type SignupReq struct {
	Account  string  `json:"account" binding:"required" example:"meowmeow789"`
	Password string  `json:"password" binding:"required" example:"meowmoew22"`
	Name     string  `json:"name" example:"Antony Cho"`
	Email    string  `json:"email" example:"meowantony@gmail.com"`
	Gender   string  `json:"gender" example:"male"`
	Phone    string  `json:"phone" example:"0912345678"`
	Birthday string  `json:"birthday" example:"2006/01/02"`
	Height   float32 `json:"height" example:"188.87"`
	Weight   float32 `json:"weight" example:"69.69"`
}

type UpdateUserInfoReq struct {
	Account        string    `json:"account" example:"meowmeow123"` // use to identify user
	Name           string    `json:"name" example:"testMeowClient"`
	Email          string    `json:"email"  binding:"email" example:"meowtestclient@gmail.com"`
	Gender         string    `json:"gender" example:"male"`
	Phone          string    `json:"phone" example:"0919886886"`
	Year           int       `json:"year" example:"2001"`
	Month          int       `json:"month" example:"5"`
	Day            int       `json:"day" example:"29"`
	Weight         float64   `json:"weight" example:"69.69"`
	Height         float64   `json:"height" example:"180.13"`
	PayType        string    `json:"pay_type" example:"visa"`
	PaymentAccount string    `json:"payment_plan" example:"1234123412341234"`
	Plan           PlanLevel `json:"plan" example:"normal"`
}
