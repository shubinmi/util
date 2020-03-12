package exec

import (
	"errors"
	"github.com/shubinmi/util/errs"
	"strings"
	"sync/atomic"
	"testing"
)

func TestAllNotNilArg(t *testing.T) {
	type args struct {
		f  o
		ps []interface{}
	}
	var n int32
	var e error

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "skip nil check",
			args: args{
				f: func(i interface{}) error {
					if i == nil {
						return errors.New("get nil")
					}
					atomic.AddInt32(&n, 1)
					return nil
				},
				ps: []interface{}{nil, &n, &n, nil, e, errors.New("")},
			},
			wantErr: false,
		},
		{
			name: "react on err",
			args: args{
				f: func(i interface{}) error {
					if e, ok := i.(error); ok {
						return e
					}
					return nil
				},
				ps: []interface{}{nil, &n, errors.New("1"), &n, errors.New("2"), nil},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AllNotNilArg(tt.args.f, tt.args.ps...)
			e = errs.Merge(e, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllNotNilArg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if n != 3 {
		t.Errorf("AllNotNilArg() must exec 3 funcs but current number %d", n)
	}
	if !strings.Contains(e.Error(), "1") || !strings.Contains(e.Error(), "2") {
		t.Errorf("AllNotNilArg() must have error with 1 and 2 in text but got %s", e.Error())
	}
}

func TestAnySuccess(t *testing.T) {
	type args struct {
		fs []func() error
	}
	var n int32
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				fs: []func() error{
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("1")
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("2")
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return nil
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("3")
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				fs: []func() error{
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("1")
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("2")
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("3")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AnySuccess(tt.args.fs...); (err != nil) != tt.wantErr {
				t.Errorf("AnySuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if n != 7 {
		t.Errorf("AnySuccess() must exec 7 funcs but current number %d", n)
	}
}

func TestUntilError(t *testing.T) {
	type args struct {
		fs []func() error
	}
	var n int32
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				fs: []func() error{
					func() error {
						atomic.AddInt32(&n, 1)
						return nil
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				fs: []func() error{
					func() error {
						atomic.AddInt32(&n, 1)
						return nil
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return nil
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("3")
					},
					func() error {
						atomic.AddInt32(&n, 1)
						return errors.New("4")
					},
				},
			},
			wantErr: true,
		},
	}
	var e error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UntilError(tt.args.fs...)
			e = errs.Merge(e, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("UntilError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if n != 5 {
		t.Errorf("UntilError() must exec 5 funcs but current number %d", n)
	}
	if !strings.Contains(e.Error(), "3") {
		t.Errorf("UntilError() must have error with 3 in text but got %s", e.Error())
	}
}

func TestUntilSuccess(t *testing.T) {
	type args struct {
		fs []func() bool
	}
	var n int32
	tests := []struct {
		name string
		args args
	}{
		{
			name: "not all",
			args: args{
				fs: []func() bool{
					func() bool {
						atomic.AddInt32(&n, 1)
						return false
					},
					func() bool {
						atomic.AddInt32(&n, 1)
						return true
					},
					func() bool {
						atomic.AddInt32(&n, 1)
						return true
					},
				},
			},
		},
		{
			name: "all",
			args: args{
				fs: []func() bool{
					func() bool {
						atomic.AddInt32(&n, 1)
						return false
					},
					func() bool {
						atomic.AddInt32(&n, 1)
						return false
					},
					func() bool {
						atomic.AddInt32(&n, 1)
						return false
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UntilSuccess(tt.args.fs...)
		})
	}

	if n != 5 {
		t.Errorf("UntilSuccess() must exec 5 funcs but current number %d", n)
	}
}
