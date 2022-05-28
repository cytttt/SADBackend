package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	MACHINE_TREADMILL   MachineName = "treadmill"   // cardio
	MACHINE_SPINNER     MachineName = "spinner"     // cardio
	MACHINE_POWER_TOWER MachineName = "power tower" // back
	MACHINE_PULLDOWN    MachineName = "pulldown"    // back
	MACHINE_LEG_PRESS   MachineName = "leg press"   // leg
	MACHINE_HACK_SQUAT  MachineName = "hack squat"  // leg
	MACHINE_AB_BENCH    MachineName = "ab bench"    // abs
	MACHINE_AB_COASTER  MachineName = "ab coaster"  // abs
	MACHINE_DONKEY_KICK MachineName = "donkey kick" // hips
	MACHINE_GHD         MachineName = "ham raise"   // hips
	MACHINE_ROWER       MachineName = "rower"       // arm
	MACHINE_BICEP_VEST  MachineName = "bicep vest"  // arm
	MACHINE_PEC_DECK    MachineName = "pec deck"    // chest
	MACHINE_CHEST_PRESS MachineName = "chest press" // chest
)

type Machine struct {
	ID              primitive.ObjectID `bson:"_id"`
	MachineID       string             `bson:"machine_id"`
	Name            string             `bson:"name"`
	Category        PartCategory       `bson:"category"`
	Gym             string             `bson:"gym_id"`
	WaitingPPL      int                `bson:"waiting_ppl"`
	BoughtAt        time.Time          `bson:"bought_at"`
	MaintainedAt    time.Time          `bson:"maintained_at"`
	ReservationOnly bool               `bson:"reservation_only"`
}

type MachineStatusResp struct {
	MachineID  string       `json:"machine_id"`
	Name       string       `json:"name"`
	WaitingPPL int          `json:"waiting_ppl"`
	Category   PartCategory `json:"category"`
}

type MachineStatusWrapperResp struct {
	Category PartCategory        `json:"category"`
	Machines []MachineStatusResp `json:"machines"`
}
