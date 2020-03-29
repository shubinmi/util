package task

import (
	"context"
	"sync"
	"time"

	"github.com/shubinmi/util/concr"
	"github.com/shubinmi/util/errs"

	uuid "github.com/satori/go.uuid"
)

const (
	WorkerClosed = iota
)

type Worker interface {
	Schedule(exec func() error, timeAfter time.Duration) (taskID string, err error)
	Status(taskID string) (status uint8, err error)
	Close()
}

type worker struct {
	q      chan *ticket
	ctx    context.Context
	cancel context.CancelFunc
	sm     StatusManager
	closed bool
	mtx    *sync.Mutex
}

func NewWorker(ctx context.Context, maxGoroutines int, sm StatusManager) Worker {
	if sm == nil {
		sm = NewImMemStatusManager()
	}
	myCtx, cancel := context.WithCancel(ctx)
	w := &worker{
		ctx:    myCtx,
		q:      make(chan *ticket, 1),
		cancel: cancel,
		sm:     sm,
		mtx:    &sync.Mutex{},
	}
	go w.run(maxGoroutines)
	go func() {
		<-myCtx.Done()
		w.Close()
	}()
	return w
}

func (w *worker) Schedule(exec func() error, timeAfter time.Duration) (taskID string, err error) {
	if w.closed {
		return "", errs.WithState(WorkerClosed, "closed")
	}
	taskID = uuid.NewV4().String()
	if err = w.sm.Save(taskID, Wait.uitn8()); err != nil {
		return
	}
	go func() {
		defer func() {
			if e := recover(); e != nil {
				return
			}
		}()
		select {
		case <-w.ctx.Done():
			return
		case <-time.After(timeAfter):
			w.q <- &ticket{
				id:   taskID,
				exec: exec,
			}
		}
	}()
	return
}

func (w *worker) Status(taskID string) (status uint8, err error) {
	if w.closed {
		return Unknown.uitn8(), errs.WithState(WorkerClosed, "closed")
	}
	return w.sm.Get(taskID)
}

func (w *worker) Close() {
	w.mtx.Lock()
	if w.closed {
		return
	}
	w.closed = true
	w.mtx.Unlock()
	w.cancel()
	close(w.q)
}

func (w *worker) run(max int) {
	addPerTime := 1
	cc := concr.NewControl(max)
	for t := range w.q {
		cc.Add(addPerTime)
		go func() {
			_ = w.sm.Save(t.id, InProgress.uitn8())
			if e := t.exec(); e != nil {
				_ = w.sm.Error(t.id, e)
			}
			_ = w.sm.Save(t.id, Done.uitn8())
			cc.Done()
		}()
		cc.WaitOnMax()
	}
}
