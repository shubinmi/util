package errs

import (
	"reflect"
	"testing"
)

func TestWithState(t *testing.T) {
	type args struct {
		state uint8
		msg   string
	}
	tests := []struct {
		name string
		args args
		want *StateErr
	}{
		{
			name: "success",
			args: args{
				state: 1,
				msg:   "1",
			},
			want: &StateErr{
				s:   1,
				msg: "1",
			},
		},
	}
	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			got := WithState(tt.args.state, tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithState() = %v, want %v", got, tt.want)
			}
			if !InState(got, tt.want.s) {
				t.Errorf("WithState() = %v, InState() failed with %v", got, tt.want.s)
			}
		})
	}
}
