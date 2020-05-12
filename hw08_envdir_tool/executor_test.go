package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	// test empty cmd
	assert.Panics(t, func() { RunCmd(nil, nil) })
	assert.Panics(t, func() { RunCmd([]string{}, nil) })

	// test non existent cmd
	assert.Error(t, errors.New("executable file not found in $PATH"), func() { RunCmd([]string{"nonexistent_cmd"}, nil) })

	// test cmd check env
	testDir := getTestDir()
	fileName := filepath.Join(testDir, "test")
	file, _ := os.Create(fileName)

	old := os.Stdout
	os.Stdout = file
	res := RunCmd([]string{"env"}, Environment{
		"ENV1": "one",
	})
	os.Stdout = old

	assert.Zero(t, res)

	data, _ := ioutil.ReadFile(fileName)
	assert.Equal(t, "ENV1=one\n", string(data))

	if err := os.Remove(fileName); err != nil {
		fmt.Println(err)
	}
	if err := os.Remove(testDir); err != nil {
		fmt.Println(err)
	}
}
