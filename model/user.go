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
	PLAN_Normal  PlanLevel = "normal"
	PLAN_Premium PlanLevel = "premium"
)
