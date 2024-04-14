package main

import (
	"os"
)

func main() {
	envDirPath := os.Args[1]
	executableCmdWithArgs := os.Args[2:]

	dirEnv, err := ReadDir(envDirPath)
	if err != nil {
		panic(err)
	}
	env := ReadCurrentEnvironemnt()

	for key, value := range dirEnv {
		env[key] = value
	}

	exitCode := RunCmd(executableCmdWithArgs, env)
	os.Exit(exitCode)
}
