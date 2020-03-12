package exec

import (
	"reflect"
	"sync"

	"github.com/shubinmi/util/errs"
)

type o func(interface{}) error

func AllNotNilArg(f o, ps ...interface{}) (e error) {
	ech := make(chan error)
	done := make(chan struct{})
	go func() {
		for err := range ech {
			e = errs.Merge(e, err)
		}
		close(done)
	}()
	wg := sync.WaitGroup{}
	for _, p := range ps {
		if reflect.ValueOf(p).Kind() == reflect.Invalid || reflect.ValueOf(p).IsNil() {
			continue
		}
		wg.Add(1)
		go func(p interface{}) {
			defer wg.Done()
			err := f(p)
			if err != nil {
				ech <- err
			}
		}(p)
	}
	wg.Wait()
	close(ech)
	<-done
	return
}

func UntilSuccess(fs ...func() bool) {
	for _, f := range fs {
		if f() {
			return
		}
	}
}

func UntilError(fs ...func() error) error {
	for _, f := range fs {
		if e := f(); e != nil {
			return e
		}
	}
	return nil
}

func AnySuccess(fs ...func() error) (err error) {
	wg := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	errors := make(chan error, 1)

	wg.Add(1)
	go func() {
		for _, f := range fs {
			wg.Add(1)
			go func(f func() error) {
				errors <- f()
				wg.Done()
			}(f)
		}
		wg.Done()
	}()

	wg2.Add(1)
	go func() {
		for e := range errors {
			if e == nil {
				err = nil
				break
			}
			err = errs.Merge(err, e)
		}
		wg2.Done()
	}()
	wg.Wait()
	close(errors)
	wg2.Wait()
	return err
}
