package errs

import (
	"errors"
	"testing"
)

func TestWithCode(t *testing.T) {
	type args struct {
		e    error
		code Code
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "access deny",
			args: args{
				e:    errors.New("1"),
				code: AccessDeny,
			},
			want: "[ecode:1]: 1",
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			err := WithCode(tt.args.e, tt.args.code)
			if err == nil {
				t.Errorf("WithCode() error = nil, wantErr %v", tt.want)
				return
			}
			if err.Error() != tt.want {
				t.Errorf("WithCode() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestExecCode(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want Code
	}{
		{
			name: "access deny",
			args: args{
				e: errors.New("access denied: [ecode:1]: action restricted for roles: []: " +
					"action restricted for roles: []"),
			},
			want: AccessDeny,
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			if got := ExecCode(tt.args.e); got != tt.want {
				t.Errorf("ExecCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCutCode(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "access deny",
			args: args{
				e: WithCode(errors.New("1"), AccessDeny),
			},
			want: "1",
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			err := CutCode(tt.args.e)
			if err == nil {
				t.Errorf("CutCode() error = nil, wantErr %v", tt.want)
				return
			}
			if err.Error() != tt.want {
				t.Errorf("CutCode() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
