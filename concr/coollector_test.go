package concr

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestErrCollector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ec := NewErrCollector(ctx)

	if ec.closed {
		t.Fatal("ec cannot be closed")
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			ec.Handle(nil)
			wg.Done()
		}()
	}
	wg.Wait()
	if e := ec.Total(); e != nil {
		t.Fatal("unexpected err", e)
	}
	if !ec.closed {
		t.Fatal("ec should be closed")
	}

	wg = sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(v int) {
			ec.Handle(fmt.Errorf(" %d", v))
			wg.Done()
		}(i)
	}
	wg.Wait()
	if ec.closed {
		t.Fatal("ec should be restarted")
	}
	e := ec.Total()
	if e == nil {
		t.Fatal("unexpected nil result")
	}
	if !ec.closed {
		t.Fatal("ec should be closed again")
	}
	errRes := e.Error() + ":"
	for i := 0; i < 10000; i++ {
		if !strings.Contains(errRes, fmt.Sprintf(" %d:", i)) {
			t.Fatal("err should contain: ", i, " got: ", e.Error())
		}
	}

	ec.Handle(nil)
	if ec.closed {
		t.Fatal("ec should be restarted again")
	}
	cancel()
	time.Sleep(10 * time.Millisecond)
	if !ec.closed {
		t.Fatal("ec should be closed on context")
	}
}
