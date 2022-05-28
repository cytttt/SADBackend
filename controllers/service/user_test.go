package service

import (
	"SADBackend/model"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessSignupInfo(t *testing.T) {
	tc1pwd := "a1b2c3"
	tc1pwdHashed := fmt.Sprintf("%x", sha256.Sum256([]byte(tc1pwd)))
	tests := []struct {
		name    string
		args    model.SignupReq
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: test password hashed result",
			args: model.SignupReq{
				Birthday: "2020/02/02",
				Password: tc1pwd,
			},
			want:    tc1pwdHashed,
			wantErr: false,
		},
		{
			name: "test case 1: missing birthday",
			args: model.SignupReq{
				Password: tc1pwd,
			},
			want:    tc1pwdHashed,
			wantErr: true,
			errMsg:  "parsing time \"\" as \"2006/01/02\": cannot parse \"\" as \"2006\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PreprocessSignupInfo(tt.args)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
			assert.Equal(t, tt.want, got.Password)
		})
	}
}

func TestVerifyPwd(t *testing.T) {
	tc1pwd := "a1b2c3"
	tc1pwdHashed := fmt.Sprintf("%x", sha256.Sum256([]byte(tc1pwd)))
	type args struct {
		req string
		ref string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errMsg  string
	}{
		{
			name: "test case 1: happy case",
			args: args{
				req: tc1pwd,
				ref: tc1pwdHashed,
			},
			wantErr: false,
		},
		{
			name: "test case 2: failure case",
			args: args{
				req: "meowmeow123",
				ref: tc1pwdHashed,
			},
			wantErr: true,
			errMsg:  "Incorrect password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPwd(tt.args.req, tt.args.ref)
			assert.Equal(t, tt.wantErr, (err != nil))
			if err != nil {
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}
		})
	}
}
