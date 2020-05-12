package main

import (
	"log"
	"os"
)

const minArgs = 3

func main() {
	args := os.Args
	if len(args) < minArgs {
		log.Fatal("Usage: go-envdir /path/to/evndir command arg1 arg2...")
	}
	dir, cmd := args[1], args[2:]

	env, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	code := RunCmd(cmd, env)
	os.Exit(code)
}
