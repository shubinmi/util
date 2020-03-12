package errs

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type Code byte

const (
	Unknown Code = iota
	AccessDeny
	NotFound
	BadRequest
	Timeout
	Internal
)

func (c Code) Alias() string {
	a, ok := map[Code]string{
		AccessDeny: "AccessDeny",
		NotFound:   "NothingFound",
		BadRequest: "BadRequest",
		Timeout:    "Timeout",
		Internal:   "Internal",
	}[c]
	if !ok {
		return "Unknown"
	}
	return a
}

func ToCode(str string) Code {
	v, e := strconv.ParseUint(str, 10, 8)
	if e != nil {
		return Unknown
	}
	c := byte(v)
	if uint64(c) != v {
		return Unknown
	}
	return Code(c)
}

func WithCode(e error, code Code) error {
	return errors.Wrapf(e, "[ecode:%v]", code)
}

func ExecCode(e error) Code {
	re := regexp.MustCompile(`\[ecode:.*]`)
	eCode := re.FindString(e.Error())
	if eCode == "" {
		return Unknown
	}
	eCode = strings.ReplaceAll(eCode, "[ecode:", "")
	eCode = strings.ReplaceAll(eCode, "]", "")
	return ToCode(eCode)
}

func CutCode(e error) error {
	re := regexp.MustCompile(`\[ecode:.*]`)
	return errors.New(re.ReplaceAllString(e.Error(), ""))
}

func SplitCode(e error) (Code, error) {
	return ExecCode(e), CutCode(e)
}
