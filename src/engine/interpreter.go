package engine

import (
	"reflect"
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
	Queue    []Command
	ExitFlag bool
	ExitChan chan bool
}

func (el *EventLoop) Start() {
	el.ExitFlag = false
	el.ExitChan = make(chan bool)
	go el.Execute()
}

func (el *EventLoop) Execute() {
	for !el.ExitFlag || len(el.Queue) > 0 {
		if len(el.Queue) > 0 {
			var cmd Command
			cmd, el.Queue = el.Queue[0], el.Queue[1:]
			cmd.Execute(el)
		}
	}
	el.ExitChan <- true
}

func (el *EventLoop) Post(cmd Command) {
	el.Queue = append(el.Queue, cmd)
}

func (el *EventLoop) AwaitFinish() {
	el.ExitFlag = true
	<-el.ExitChan
}
