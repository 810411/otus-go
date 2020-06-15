package main

import (
	"errors"
	"log"
	"os"

	"github.com/810411/otus-go/hw09_generator_of_validators/go-validate/internal"
)

const osArgs = 2

var ErrFileNotRegular = errors.New("file not regular")

func main() {
	args := os.Args
	if len(args) < osArgs {
		log.Fatal("Usage: go-validate /path/to/models.go(path to file with structs description)")
	}
	filepath := args[1]
	if err := fileInPath(filepath); err != nil {
		log.Fatal(err)
	}

	packName, stForValidArr, err := internal.Parse(filepath)
	if err != nil {
		log.Fatal(err)
	}

	if err = internal.Generate(filepath, packName, stForValidArr); err != nil {
		log.Fatal(err)
	}
}

func fileInPath(path string) error {
	stat, err := os.Stat(path)

	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		return ErrFileNotRegular
	}

	return nil
}
