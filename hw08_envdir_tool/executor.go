package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	name, args := cmd[0], cmd[1:]

	proc := exec.Command(name, args...)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin

	if env != nil {
		envs := make([]string, 0, len(env))
		for k, v := range env {
			envs = append(envs, k+"="+v)
		}

		proc.Env = envs
	}

	if err := proc.Run(); err != nil {
		log.Fatal(err)
	}

	return 0
}
