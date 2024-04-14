package main

import (
	"errors"
	"os"
	"os/exec"
)

var (
	ErrCommandRequired = errors.New("command argument not passed")
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		returnCode = 1
		return
	}

	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Env = createEnv(env)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Start(); err != nil {
		returnCode = 1
		return
	}

	if err := execCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			returnCode = exitErr.ExitCode()
			return
		} else {
			returnCode = 1
			return
		}
	}

	return
}

func createEnv(env Environment) []string {
	envStrings := []string{}
	for key, val := range env {
		if val.NeedRemove {
			continue
		}
		envStrings = append(envStrings, key+"="+val.Value)
	}
	return envStrings
}
