package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Machine struct {
	ID           primitive.ObjectID `bson:"_id"`
	MachineID    string             `bson:"machine_id"`
	Name         string             `bson:"name"`
	Category     PartCategory       `bson:"category"`
	Gym          string             `bson:"gym_id"`
	WaitingPPL   int                `bson:"waiting_ppl"`
	BoughtAt     time.Time          `bson:"bought_at"`
	MaintainedAt time.Time          `bson:"maintained_at"`
}

type PartCategory string

const (
	PART_BACK   PartCategory = "back"
	PART_CHEST  PartCategory = "chest"
	PART_CARDIO PartCategory = "cardio"
	PART_ABS    PartCategory = "abs" // abdomen
	PART_LEG    PartCategory = "leg"
	PART_ARM    PartCategory = "arm"
	PART_HIPS   PartCategory = "hips"
)

type MachineName string

const (
	MACHINE_TREADMILL   = "treadmill"   // cardio
	MACHINE_SPINNER     = "spinner"     // cardio
	MACHINE_POWER_TOWER = "power tower" // back
	MACHINE_PULLDOWN    = "pulldown"    // back
	MACHINE_LEG_PRESS   = "leg press"   // leg
	MACHINE_HACK_SQUAT  = "hack squat"  // leg
	MACHINE_AB_BENCH    = "ab bench"    // abs
	MACHINE_AB_COASTER  = "ab coaster"  // abs
	MACHINE_DONKEY_KICK = "donkey kick" // hips
	MACHINE_GHD         = "ham raise"   // hips
	MACHINE_ROWER       = "rower"       // arm
	MACHINE_BICEP_VEST  = "bicep vest"  // arm
	MACHINE_PEC_DECK    = "pec deck"    // chest
	MACHINE_CHEST_PRESS = "chest press" // chest
)
