package concr

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestControl(t *testing.T) {
	cc := NewControl(3)
	if cc.inWait {
		t.Error("wrong time for inWait true")
	}
	cc.Add(2)
	if cc.inWait {
		t.Error("wrong time for inWait true")
	}
	cc.Add(1)
	if !cc.inWait {
		t.Error("wrong time for inWait false")
	}
	var i int32
	go func() {
		cc.WaitOnMax()
		atomic.AddInt32(&i, 1)
	}()
	time.Sleep(100 * time.Millisecond)
	if i > 0 {
		t.Error("wrong WaitOnMax logic when should wait")
	}
	cc.Done()
	go func() {
		cc.WaitOnMax()
		atomic.AddInt32(&i, 1)
	}()
	time.Sleep(100 * time.Millisecond)
	if i != 2 {
		t.Error("wrong WaitOnMax logic when should be done")
	}
}

func TestControl_WaitAll(t *testing.T) {
	cc := NewControl(50)
	var s int64
	for i := 0; i < 1000; i++ {
		cc.Add(1)
		go func() {
			atomic.AddInt64(&s, 1)
			cc.Done()
		}()
		cc.WaitOnMax()
		if cc.currentNum > 50 {
			t.Error("concurrent goroutines more then limit")
		}
	}
	cc.WaitAll()
	if s != 1000 {
		t.Errorf("resutl should be 100 bu got %d", s)
	}
}
