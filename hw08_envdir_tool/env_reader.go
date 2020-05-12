package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	var result Environment

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		b, err := ioutil.ReadFile(filepath.Join(dir, filename))
		if err != nil {
			return nil, err
		}
		if len(b) == 0 {
			delete(result, filename)
			continue
		}
		b = bytes.ReplaceAll(bytes.Split(b, []byte("\n"))[0], []byte("\x00"), []byte("\n"))

		value := strings.TrimRight(string(b), " \t\n")
		if strings.ContainsRune(value, '=') {
			continue
		}

		if result == nil {
			result = make(Environment)
		}
		result[filename] = value
	}

	return result, err
}
