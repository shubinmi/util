package errs

import "github.com/pkg/errors"

func Merge(ers ...error) (r error) {
	for _, e := range ers {
		if e == nil {
			continue
		}
		if r == nil {
			r = e
			continue
		}
		r = errors.Wrap(r, e.Error())
	}
	return
}

func WrapForce(e error, msg string) error {
	if e == nil {
		return errors.New(msg)
	}
	return errors.Wrap(e, msg)
}
