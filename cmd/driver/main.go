package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()

	for _, filename := range args {
		err := runCase(filename)
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		}
	}
}
