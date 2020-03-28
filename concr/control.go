package concr

import (
	"sync"
	"time"
)

type Control struct {
	currentNum int
	maxNum     int
	wait       chan struct{}
	inWait     bool
	mtx        *sync.Mutex
}

func NewControl(maxGoroutines int) *Control {
	ch := make(chan struct{})
	close(ch)
	return &Control{
		maxNum: maxGoroutines,
		wait:   ch,
		mtx:    &sync.Mutex{},
	}
}

func (mc *Control) Add(num int) {
	mc.mtx.Lock()
	mc.currentNum += num
	mc.mtx.Unlock()
	if mc.currentNum >= mc.maxNum {
		mc.inWait = true
		mc.wait = make(chan struct{})
	}
}

func (mc *Control) Done() {
	mc.mtx.Lock()
	mc.currentNum--
	mc.mtx.Unlock()
	if mc.inWait && mc.currentNum < mc.maxNum {
		defer func() {
			if e := recover(); e != nil {
				return
			}
		}()
		mc.inWait = false
		close(mc.wait)
	}
}

func (mc *Control) WaitOnMax() {
	<-mc.wait
}

func (mc *Control) WaitAll() {
	sleep := 10 * time.Nanosecond
	for mc.currentNum > 0 {
		time.Sleep(sleep)
	}
}
