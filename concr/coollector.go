package concr

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/shubinmi/util/errs"
)

type errState uint8

const (
	collected errState = iota + 1
)

type ErrCollector struct {
	ch        chan error
	err       error
	ctx       context.Context
	closed    bool
	collected chan struct{}
	mtx       *sync.Mutex
}

func NewErrCollector(ctx context.Context) *ErrCollector {
	ec := &ErrCollector{
		ch:        make(chan error, 1),
		ctx:       ctx,
		mtx:       &sync.Mutex{},
		collected: make(chan struct{}),
	}
	go ec.run()
	return ec
}

func (ec *ErrCollector) run() {
	for {
		select {
		case <-ec.ctx.Done():
			_ = ec.Total()
			return
		case err, ok := <-ec.ch:
			if !ok {
				return
			}
			if err == nil {
				continue
			}
			if errs.InState(err, collected) {
				ec.collected <- struct{}{}
				return
			}
			ec.err = errs.Merge(ec.err, err)
		}
	}
}

func (ec *ErrCollector) Handle(e error) {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			return
		}
	}()
	ec.mtx.Lock()
	defer ec.mtx.Unlock()
	if ec.closed {
		ec.closed = false
		ec.err = nil
		go ec.run()
		waitReadChBeforeWrite := 100 * time.Nanosecond
		time.Sleep(waitReadChBeforeWrite)
	}
	ec.ch <- e
}

func (ec *ErrCollector) Total() error {
	ec.mtx.Lock()
	defer ec.mtx.Unlock()
	if !ec.closed {
		ec.closed = true
		ec.ch <- errs.WithState(collected, "")
		<-ec.collected
	}
	return ec.err
}
