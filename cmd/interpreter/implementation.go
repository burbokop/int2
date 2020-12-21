package main

import (
	"fmt"

	engine "github.com/burbokop/int2/src/engine"
)

//print <str>
type PrintCommand struct {
	arg string
}

func (p *PrintCommand) Init(args []string) {
	if len(args) > 1 {
		p.arg = args[1]
	}
}

func (p *PrintCommand) Execute(loop engine.Handler) {
	fmt.Println(p.arg)
}

//cat <arg1> <arg2>
type CatCommand struct {
	str0 string
	str1 string
}

func (p *CatCommand) Init(args []string) {
	if len(args) > 2 {
		p.str0 = args[1]
		p.str1 = args[2]
	}
}

func (p *CatCommand) Execute(loop engine.Handler) {
	loop.Post(&PrintCommand{arg: p.str0 + p.str1})
}
