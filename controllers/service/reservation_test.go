package service

import (
	"SADBackend/model"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComputeStartTime(t *testing.T) {
	tc1, _ := string2Time("2022-02-04 18:30", "2006-01-02 15:04")
	type args struct {
		ymd string
		hm  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: happy case",
			args: args{
				ymd: "2022-02-04",
				hm:  "18:30",
			},
			want:    *tc1,
			wantErr: false,
		},
		{
			name: "test case 2: wrong format",
			args: args{
				ymd: "2022/02/04",
				hm:  "18:30",
			},
			wantErr: true,
			errMsg:  "parsing time \"2022/02/04\" as \"2006-01-02\": cannot parse \"/02/04\" as \"-\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComputeStartTime(tt.args.ymd, tt.args.hm)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetTimeRangeBound(t *testing.T) {
	type args struct {
		tr   string
		date string
	}
	type res struct {
		lh int
		lm int
		uh int
		um int
	}
	tests := []struct {
		name    string
		args    args
		want    res
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: morning",
			args: args{
				tr:   "morning",
				date: "2022-02-20",
			},
			wantErr: false,
			want: res{
				lh: 6,
				lm: 0,
				uh: 11,
				um: 30,
			},
		},
		{
			name: "test case 2: afternoon",
			args: args{
				tr:   "afternoon",
				date: "2022-02-20",
			},
			wantErr: false,
			want: res{
				lh: 12,
				lm: 0,
				uh: 17,
				um: 30,
			},
		},
		{
			name: "test case 3: evening",
			args: args{
				tr:   "evening",
				date: "2022-02-20",
			},
			wantErr: false,
			want: res{
				lh: 18,
				lm: 0,
				uh: 20,
				um: 30,
			},
		},
		{
			name: "test case 4: midnight",
			args: args{
				tr:   "midnight",
				date: "2022-02-20",
			},
			wantErr: false,
			want: res{
				lh: 21,
				lm: 0,
				uh: 23,
				um: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll, uu, err := GetTimeRangeBound(tt.args.tr, tt.args.date)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			llh, llm, _ := ll.Clock()
			uuh, uum, _ := uu.Clock()
			assert.Equal(t, tt.want.lh, llh)
			assert.Equal(t, tt.want.lm, llm)
			assert.Equal(t, tt.want.uh, uuh)
			assert.Equal(t, tt.want.um, uum)
		})
	}
}

func TestFindKAvailabeTime(t *testing.T) {
	rand.Seed(42)
	loc := time.FixedZone("Asia/Taipei", int((8 * time.Hour).Seconds()))
	lb := time.Date(2022, 9, 10, 12, 0, 0, 0, loc)
	ub := time.Date(2022, 9, 10, 15, 0, 0, 0, loc)
	k := 5
	type args struct {
		userID string
		res    []model.Reservation
		mID    []string
		lb     time.Time
		ub     time.Time
		k      int
	}
	tests := []struct {
		name    string
		args    args
		want    []model.GetAvailableTimeResp
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: happy case",
			args: args{
				userID: "test-user-1",
				mID: []string{
					"test-machine-1",
					"test-machine-2",
				},
				lb: lb,
				ub: ub,
				k:  k,
			},
			wantErr: false,
			want: []model.GetAvailableTimeResp{
				{Start: "12:30", End: "13:00", MachineID: "test-machine-2"},
				{Start: "13:30", End: "14:00", MachineID: "test-machine-1"},
				{Start: "14:00", End: "14:30", MachineID: "test-machine-2"},
				{Start: "14:30", End: "15:00", MachineID: "test-machine-2"},
				{Start: "15:00", End: "15:30", MachineID: "test-machine-2"},
			},
		},
		{
			name: "test case 1: very large k",
			args: args{
				userID: "test-user-1",
				mID: []string{
					"test-machine-1",
					"test-machine-2",
				},
				lb: lb,
				ub: ub,
				k:  1000,
			},
			wantErr: false,
			want: []model.GetAvailableTimeResp{
				{Start: "12:00", End: "12:30", MachineID: "test-machine-2"},
				{Start: "12:30", End: "13:00", MachineID: "test-machine-2"},
				{Start: "13:00", End: "13:30", MachineID: "test-machine-1"},
				{Start: "13:30", End: "14:00", MachineID: "test-machine-1"},
				{Start: "14:00", End: "14:30", MachineID: "test-machine-2"},
				{Start: "14:30", End: "15:00", MachineID: "test-machine-2"},
				{Start: "15:00", End: "15:30", MachineID: "test-machine-2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(42)
			got, err := FindKAvailabeTime(tt.args.userID, tt.args.res, tt.args.mID, tt.args.lb, tt.args.ub, tt.args.k)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
