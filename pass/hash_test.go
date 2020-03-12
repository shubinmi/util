package pass

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashVerify(t *testing.T) {
	tests := map[string]struct {
		pass   string
		suffix string
		errHas string
	}{
		"match": {pass: "pass123", errHas: ""},
		"miss":  {pass: "#asdas", suffix: "weq", errHas: bcrypt.ErrMismatchedHashAndPassword.Error()},
	}
	for name, tt := range tests {
		test := tt
		t.Run(name, func(t *testing.T) {
			hash, err := Hash(test.pass)
			if err != nil {
				if test.errHas == "" {
					t.Fatal("got unexpected err", err)
				}
				if !strings.Contains(err.Error(), test.errHas) {
					t.Fatal("err doesn't match", err.Error(), test.errHas)
				}
				return
			}
			err = Verify(hash+test.suffix, test.pass)
			if err != nil {
				if test.errHas == "" {
					t.Fatal("got unexpected err", err)
				}
				if !strings.Contains(err.Error(), test.errHas) {
					t.Fatal("err doesn't match", err.Error(), test.errHas)
				}
				return
			}
		})
	}
}
