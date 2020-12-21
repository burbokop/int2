package main

import (
	"fmt"
	"strings"

	engine "github.com/burbokop/simple_interpreter/src/engine"
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

//reverse <str>
type ReverseCommand struct {
	arg       string
	deferMode bool
}

func (p *ReverseCommand) Init(args []string) {
	if len(args) > 1 {
		p.arg = args[1]
		if len(args) > 2 {
			if args[2] == "defer" {
				p.deferMode = true
			}
		}
	}
}

func (p *ReverseCommand) Execute(loop engine.Handler) {
	var reverce = func(s string) string {
		rns := []rune(s)
		for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
			rns[i], rns[j] = rns[j], rns[i]
		}
		return string(rns)
	}
	if p.deferMode {
		loop.Post(&PrintCommand{arg: reverce(p.arg)})
	} else {
		fmt.Println(reverce(p.arg))
	}
}

//split <str> <sep>
type SplitCommand struct {
	str       string
	sep       string
	deferMode bool
}

func (p *SplitCommand) Init(args []string) {
	if len(args) > 2 {
		p.str = args[1]
		p.sep = args[2]
		if len(args) > 3 {
			if args[3] == "defer" {
				p.deferMode = true
			}
		}
	}
}

func (p *SplitCommand) Execute(loop engine.Handler) {
	var r = strings.Split(p.str, p.sep)
	if p.deferMode {
		for _, e := range r {
			loop.Post(&PrintCommand{arg: e})
		}
	} else {
		for _, e := range r {
			fmt.Println(e)
		}
	}
}
