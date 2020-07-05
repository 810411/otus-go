package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

const wait = 13 * time.Second

var host string

func TestMain(m *testing.M) {
	host = os.Getenv("HOST")

	log.Printf("waiting %s while starting", wait)
	time.Sleep(wait)

	opts := godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: 0,
	}

	status := godog.TestSuite{
		Name:                "integration-tests",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
