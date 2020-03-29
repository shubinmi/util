package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shubinmi/util/errs"
)

func TestWorker(t *testing.T) {
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())

	w := NewWorker(ctx1, 1, NewImMemStatusManager())
	defer w.Close()

	var i int32
	f := func() error {
		<-ctx2.Done()
		atomic.AddInt32(&i, 1)
		return nil
	}

	id1, e := w.Schedule(f, 0)
	if e != nil {
		t.Errorf("unexpected err on first call = %v", e)
	}
	// because goroutines could be run in random order
	time.Sleep(20 * time.Millisecond)
	id2, e := w.Schedule(f, 0)
	if e != nil {
		t.Errorf("unexpected err on second call = %v", e)
	}

	// First checks
	time.Sleep(50 * time.Millisecond)
	if i != 0 {
		t.Errorf("wrong behavior1, i should be = 0, but got = %d", i)
	}
	s, _ := w.Status(id1)
	if s != InProgress.uitn8() {
		t.Errorf("wrong behavior1, status of first should be = 2, but got = %d", s)
	}
	s, _ = w.Status(id2)
	if s != Wait.uitn8() {
		t.Errorf("wrong behavior1, status of second should be = 1, but got = %d", s)
	}

	cancel2()
	time.Sleep(50 * time.Millisecond)
	// Second checks
	if i != 2 {
		t.Errorf("wrong behavior2, i should be = 2, but got = %d", i)
	}
	s, e = w.Status(id1)
	if s != Done.uitn8() {
		t.Errorf("wrong behavior2, status of first should be = 3, but got = %d", s)
	}
	if e != nil {
		t.Errorf("unexpected err on first done = %v", e)
	}
	s, e = w.Status(id2)
	if s != Done.uitn8() {
		t.Errorf("wrong behavior2, status of second should be = 3, but got = %d", s)
	}
	if e != nil {
		t.Errorf("unexpected err on second done = %v", e)
	}
	time.Sleep(50 * time.Millisecond)

	cancel1()
	time.Sleep(50 * time.Millisecond)
	// Third checks
	_, e = w.Schedule(f, 0)
	if !errs.InState(e, WorkerClosed) {
		t.Error("worker should be closed")
	}
}
