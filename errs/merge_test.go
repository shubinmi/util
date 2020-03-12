package errs

import (
	"errors"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	type args struct {
		ers []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil",
			args: args{
				ers: []error{nil, nil, nil},
			},
			wantErr: false,
		},
		{
			name: "nil",
			args: args{
				ers: []error{nil, errors.New("1"), nil, errors.New("2")},
			},
			wantErr: true,
		},
	}
	var e error
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			err := Merge(tt.args.ers...)
			e = Merge(e, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	if !strings.Contains(e.Error(), "1") || !strings.Contains(e.Error(), "2") {
		t.Errorf("Merge() wait error with 1 and 2 in text but got %s", e.Error())
	}
}
