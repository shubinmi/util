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
