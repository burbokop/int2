package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"reflect"

	engine "github.com/burbokop/simple_interpreter/src/engine"
)

func main() {
	var inputPath string
	flag.StringVar(&inputPath, "i", "", "Input file path")
	flag.Parse()
	if len(inputPath) == 0 {
		return
	}

	var parser = new(engine.Parser)
	parser.AddCmdType("print", reflect.TypeOf(PrintCommand{}))
	parser.AddCmdType("reverse", reflect.TypeOf(ReverseCommand{}))
	parser.AddCmdType("split", reflect.TypeOf(SplitCommand{}))

	eventLoop := new(engine.EventLoop)
	eventLoop.Start()

	if input, err := os.Open(inputPath); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd, err := parser.Parse(commandLine) // parse the line to get an instance of Command
			if err == nil {
				eventLoop.Post(cmd)
			} else {
				log.Fatal(err)
			}
		}
	}
	eventLoop.AwaitFinish()
}
