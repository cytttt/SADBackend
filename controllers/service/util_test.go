package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString2Time(t *testing.T) {
	offset := int((8 * time.Hour).Seconds())
	loc := time.FixedZone("Asia/Taipei", offset)
	type args struct {
		timeStr string
		format  string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
		ErrMsg  string
	}{
		{
			name: "test case 1: fail, wrong time format",
			args: args{
				timeStr: "2021/7/21 10:58",
				format:  "2006/01/02 15:04",
			},
			want:    time.Date(2021, 7, 21, 10, 58, 0, 0, loc),
			wantErr: true,
			ErrMsg:  "parsing time \"2021/7/21 10:58\" as \"2006/01/02 15:04\": cannot parse \"7/21 10:58\" as \"01\"",
		},
		{
			name: "test case 2: correct",
			args: args{
				timeStr: "2021/07/21 10:58",
				format:  "2006/01/02 15:04",
			},
			want:    time.Date(2021, 7, 21, 10, 58, 0, 0, loc),
			wantErr: false,
		},
		{
			name: "test case 3: fail, exceed time value",
			args: args{
				timeStr: "2021/07/32 10:58",
				format:  "2006/01/02 15:04",
			},
			want:    time.Date(2021, 8, 1, 10, 58, 0, 0, loc),
			wantErr: true,
			ErrMsg:  "parsing time \"2021/07/32 10:58\": day out of range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := string2Time(tt.args.timeStr, tt.args.format)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.ErrMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want, *got)
		})
	}
}
