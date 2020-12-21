package engine

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/enriquebris/goconcurrentqueue"
)

type Handler interface {
	Post(cmd Command)
}

type Command interface {
	Init(args []string)
	Execute(handler Handler)
}

func CommandType() reflect.Type { return reflect.TypeOf((*Command)(nil)).Elem() }

type EventLoop struct {
	Queue        *goconcurrentqueue.FIFO
	Active       bool
	ExecFinished sync.WaitGroup
	ExecAlive    int32
}

func (el *EventLoop) Exec() {
	defer atomic.StoreInt32(&el.ExecAlive, 0)
	defer el.ExecFinished.Done()
	if el.Queue == nil {
		return
	}
	for el.Queue.GetLen() > 0 {
		var cmd, err = el.Queue.Dequeue()
		if err == nil {
			cmd.(Command).Execute(el)
		}
	}
}

func (el *EventLoop) Start() {
	el.Active = true
	if atomic.CompareAndSwapInt32(&el.ExecAlive, 0, 1) {
		el.ExecFinished.Add(1)
		go el.Exec()
	}
}

func (el *EventLoop) Post(cmd Command) {
	if el.Queue == nil {
		el.Queue = goconcurrentqueue.NewFIFO()
	}
	el.Queue.Enqueue(cmd)
	if el.Active && atomic.CompareAndSwapInt32(&el.ExecAlive, 0, 1) {
		el.ExecFinished.Add(1)
		go el.Exec()
	}
}

func (el *EventLoop) AwaitFinish() {
	el.Active = false
	el.ExecFinished.Wait()
}
