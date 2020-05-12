package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	testDir := getTestDir()

	// test incorrect dir
	incorrectDir := filepath.Join(testDir, "incorrect_dir")
	env, err := ReadDir(incorrectDir)
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
	assert.Nil(t, env)

	// test path is not dir
	file := createFile("FILE", []byte(""))
	env, err = ReadDir(file)
	assert.Error(t, err)
	assert.Nil(t, env)

	if err := os.Remove(file); err != nil {
		fmt.Println(err)
	}

	// test empty dir
	env, err = ReadDir(testDir)
	assert.NoError(t, err)
	assert.Nil(t, env)

	if err := os.Remove(testDir); err != nil {
		fmt.Println(err)
	}
}

func getTestDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(dir, "test_temp")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return dir
}

func createFile(fileName string, data []byte) string {
	testDir := getTestDir()
	filePath := filepath.Join(testDir, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		log.Panicf("failed to create: %v", err)
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		log.Panicf("failed to write: %v", err)
	}

	return filePath
}
