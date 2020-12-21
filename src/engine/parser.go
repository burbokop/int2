package engine

import (
	"fmt"
	"reflect"
	"strings"
)

type Parser struct {
	Cmds map[string]reflect.Type
}

type ImplementationError struct {
	CurrentType reflect.Type
	NeededType  reflect.Type
}

func (e *ImplementationError) Error() string {
	return "Implementation error: Type " + e.CurrentType.Name() + " must implement " + e.NeededType.Name()
}

type UnknownCmdError struct {
	CmdName string
}

func (e *UnknownCmdError) Error() string {
	return "Unknown command: " + e.CmdName
}

type CommandCreationError struct {
	Type reflect.Type
}

func (e *CommandCreationError) Error() string {
	return "Can not create command of type: " + e.Type.Name()
}

type EmptyLineError struct {
	Type reflect.Type
}

func (e *EmptyLineError) Error() string {
	return "Empty line"
}

func (parser *Parser) AddCmdType(name string, t reflect.Type) error {
	t = RemovePtr(t)
	var ptrType = reflect.PtrTo(t)
	if parser.Cmds == nil {
		parser.Cmds = make(map[string]reflect.Type)
	}
	if ptrType.Implements(CommandType()) {
		parser.Cmds[name] = t
		return nil
	} else {
		return &ImplementationError{t, CommandType()}
	}
}

func RemovePtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func (parser *Parser) Parse(str string) (Command, error) {
	var s = strings.Fields(str)
	if len(s) > 0 {
		var t, found = parser.Cmds[s[0]]
		if found {
			var ptr = reflect.New(t)
			if ptr.IsNil() {
				return nil, &CommandCreationError{t}
			} else {
				var cmd_ptr = ptr.Interface().(Command)
				cmd_ptr.Init(s)
				return cmd_ptr, nil
			}
		} else {
			return nil, &UnknownCmdError{s[0]}
		}
	}
	return nil, &EmptyLineError{}
}

func Print(cmds []Command) {
	for i, cmd := range cmds {
		fmt.Println(i, reflect.TypeOf(cmd).Elem().Name(), cmd)
	}
}
