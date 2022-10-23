package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesEnv, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "error read dir", err)
	}

	envs := make(Environment, len(filesEnv))
	for _, fileInfo := range filesEnv {
		line, err := getFirstLine(dir, fileInfo.Name())
		if err != nil {
			return nil, err
		}

		if strings.ContainsAny(line, "=") {
			return nil, fmt.Errorf("%s %s: %w", "file has forbidden chars", fileInfo.Name(), err)
		}

		envValue := strings.TrimRight(line, " ")
		if envValue == " \n" {
			envValue = ""
		} else {
			envValue = strings.TrimRight(line, "\n")
		}
		envValue = string(bytes.ReplaceAll([]byte(envValue), []byte{0x00}, []byte("\n")))

		envs[fileInfo.Name()] = EnvValue{
			Value:      envValue,
			NeedRemove: envValue == "",
		}
	}

	return envs, nil
}

func getFirstLine(dir string, fileName string) (string, error) {
	file, err := os.Open(dir + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", "error open file", fileName, err)
	}

	scanner := bufio.NewReader(file)
	line, err := scanner.ReadString('\n')
	if len(line) == 0 && err != nil {
		if errors.Is(err, io.EOF) {
			return "", nil
		}
		return "", fmt.Errorf("%s %s: %w", "error read file", fileName, err)
	}

	return line, nil
}
