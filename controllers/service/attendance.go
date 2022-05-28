package service

import (
	"SADBackend/model"
	"sort"
)

func PostprocessStatData(raw []model.StatInSecond) ([]model.CompanyStatResp, error) {
	sort.Slice(raw, func(i, j int) bool {
		return raw[i].Date < raw[j].Date
	})
	var results []model.CompanyStatResp
	for _, i := range raw {
		results = append(results, model.CompanyStatResp{
			Date:        i.Date,
			Attendance:  i.AttendanceCount,
			AvgStayTime: float32(i.AvgStaySecond) / 3600,
		})
	}
	return results, nil
}
