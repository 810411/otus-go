package main

import (
	"os"
	"os/exec"
)

const (
	ExitSuccess = 0
	ExitFailure = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return ExitFailure
	}

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
		return ExitFailure
	}

	return ExitSuccess
}
