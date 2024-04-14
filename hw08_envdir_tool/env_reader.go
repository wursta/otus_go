package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

var ErrUnsupportedFile = errors.New("unsupported file")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		envValue, err := readEnvFile(dir, file.Name())
		if err != nil {
			return nil, err
		}

		env[file.Name()] = envValue
	}

	return env, nil
}

func ReadCurrentEnvironemnt() Environment {
	env := Environment{}

	for _, envString := range os.Environ() {
		keyValue := strings.Split(envString, "=")
		env[keyValue[0]] = EnvValue{
			Value: keyValue[1],
		}
	}

	return env
}

func readEnvFile(dirPath, fileName string) (envValue EnvValue, err error) {
	if strings.ContainsAny(fileName, ".=") {
		err = ErrUnsupportedFile
		return
	}

	file, err := os.Open(dirPath + "/" + fileName)
	if err != nil {
		return
	}

	defer file.Close()

	sc := bufio.NewScanner(file)
	var line int
	var firstLineText string

	for sc.Scan() {
		line++
		if line == 1 {
			err = sc.Err()
			if err != nil {
				return
			}

			firstLineText = sc.Text()
		}
	}

	firstLineText = strings.TrimRight(firstLineText, " ")
	firstLineText = string(bytes.ReplaceAll([]byte(firstLineText), []byte{0x00}, []byte("\n")))

	envValue.Value = firstLineText
	if firstLineText == "" {
		envValue.NeedRemove = true
	}
	return
}
