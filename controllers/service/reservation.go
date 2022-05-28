package service

import (
	"SADBackend/model"
	"math/rand"
	"sort"
	"time"

	mapset "github.com/deckarep/golang-set"
)

func FindFutureKReservation(raw []model.AggrReservationRes, k int) ([]model.ReservationResp, error) {
	var results []model.ReservationResp
	for _, i := range raw {
		results = append(results, model.ReservationResp{
			MachineID:   i.MachineID,
			Category:    string(i.Machines[0].Category),
			MachineName: i.Machines[0].Name,
			GymID:       i.Gyms[0].BranchGymID,
			GymName:     i.Gyms[0].Name,
			Date:        i.StartAt,
		})
		if len(results) == k {
			break
		}
	}
	return results, nil
}

func ComputeStartTime(yyyymmdd, hhmm string) (time.Time, error) {
	startBase, err := string2Time(yyyymmdd, "2006-01-02")
	if err != nil {
		return time.Time{}, err
	}
	startHour, err := time.Parse("15:04", hhmm)
	if err != nil {
		return time.Time{}, err
	}
	start := startBase.Add(time.Duration(startHour.Hour()) * time.Hour)
	start = start.Add(time.Duration(startHour.Minute()) * time.Minute)
	return start, nil
}

func FindKAvailabeTime(userID string, res []model.Reservation, mID []string, lbd time.Time, ubd time.Time, k int) ([]model.GetAvailableTimeResp, error) {
	var tmpAvailableTime []model.GetAvailableTimeResp
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
		tmpAvailableTime = append(tmpAvailableTime, model.GetAvailableTimeResp{
			Start:     st,
			End:       ed,
			MachineID: mIDAtd[rand.Intn(len(mIDAtd))],
		})
	}
	permutation := rand.Perm(len(tmpAvailableTime))[:min(k, len(tmpAvailableTime))]
	sort.Ints(permutation)
	var AvailableTime []model.GetAvailableTimeResp
	for _, p := range permutation {
		AvailableTime = append(AvailableTime, tmpAvailableTime[p])
	}
	return AvailableTime, nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
