package task

import (
	"sync"

	"github.com/shubinmi/util/errs"
)

type Status uint8

// noinspection GoUnusedConst
const (
	Unknown Status = iota
	Wait
	InProgress
	Done
)

func (s Status) uitn8() uint8 {
	return uint8(s)
}

type StatusManager interface {
	Save(taskID string, status uint8) error
	Error(taskID string, err error) error
	Get(taskID string) (status uint8, err error)
}

type taskProgress struct {
	id     string
	status uint8
	error  error
}

type imMemStatusManager struct {
	data sync.Map
}

func NewImMemStatusManager() StatusManager {
	return &imMemStatusManager{
		data: sync.Map{},
	}
}

func (i *imMemStatusManager) Save(taskID string, status uint8) error {
	tp := taskProgress{id: taskID}
	v, ok := i.data.Load(taskID)
	if ok {
		tp = v.(taskProgress)
	}
	tp.status = status
	i.data.Store(taskID, tp)
	return nil
}

func (i *imMemStatusManager) Error(taskID string, err error) error {
	tp := taskProgress{id: taskID}
	v, ok := i.data.Load(taskID)
	if ok {
		tp = v.(taskProgress)
	}
	tp.error = err
	i.data.Store(taskID, tp)
	return nil
}

func (i *imMemStatusManager) Get(taskID string) (status uint8, err error) {
	tp := taskProgress{id: taskID}
	v, ok := i.data.Load(taskID)
	if !ok {
		return Unknown.uitn8(), errs.NothingFound{}
	}
	tp = v.(taskProgress)
	return tp.status, tp.error
}
