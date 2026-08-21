package main

import (
	"os"
	"path"

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

func runCaseFile(filename string) error {
	driverItem, err := initDriver(filename)
	if err != nil {
		return err
	}

	dirName := path.Dir(filename)
	return driver.GenericRun(driverItem, dirName)
}

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return stat.IsDir(), nil
}

func runCase(filename string) error {
	isDir, err := isDir(filename)
	if err != nil {
		return err
	}

	if !isDir {
		return runCaseFile(filename)
	}

	files, err := os.ReadDir(filename)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			fullname := path.Join(filename, file.Name())
			err := runCase(fullname)
			if err != nil {
				return err
			}
			continue
		}

		if path.Ext(file.Name()) != ".bftest" {
			continue
		}

		err := runCaseFile(path.Join(filename, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
