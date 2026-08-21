package main

import (
	"flag"
	"fmt"
)

const (
	CommandRun = "run"
)

func doRun(args []string) {
	set := flag.NewFlagSet("run", flag.ExitOnError)
	_ = set.Parse(args)

	if set.NArg() <= 0 {
		set.Usage()
		return
	}

	for _, filename := range set.Args() {
		err := runCase(filename)
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		}
	}
}

func main() {
	flag.Parse()

	if flag.NArg() <= 0 {
		fmt.Printf("Usage: %s [COMMAND]\n", "go-brainfuck")
		fmt.Printf("commands:\n")
		fmt.Printf("  run [FILE]...  run test driver\n")
		return
	}

	args := flag.Args()
	switch args[0] {
	case CommandRun:
		doRun(args[1:])

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
	}
}
