package main

import (
	"github.com/flily/go-brainfuck/context"
	"github.com/flily/go-brainfuck/driver"
)

func initDriver(filename string) (*driver.TestDriverItem, error) {
	fctx, err := context.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parser := driver.NewParser(filename, fctx)
	driverItem, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return driverItem, nil
}

func runCase(filename string) error {
	driverItem, err := initDriver(filename)
	if err != nil {
		return err
	}

	return driver.GenericRun(driverItem)
}
