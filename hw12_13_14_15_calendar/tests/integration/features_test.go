package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

type feature struct {
	resp *http.Response
}

func (f *feature) iSendRequestTo(method, endpoint string) error {
	return f.iSendRequestToWithJsonBody(method, endpoint, nil)
}

func (f *feature) iSendRequestToWithJsonBody(method, endpoint string, message *messages.PickleStepArgument_PickleDocString) error {
	var body io.Reader
	if message != nil {
		body = strings.NewReader(message.Content)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	f.resp, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (f *feature) theResponseCodeShouldBe(code int) error {
	if code != f.resp.StatusCode {
		var actual interface{}
		err := json.NewDecoder(f.resp.Body).Decode(&actual)
		if err != nil {
			return err
		}
		fmt.Println(actual)

		return fmt.Errorf("expected response code to be: %d, but actual is: %d", code, f.resp.StatusCode)
	}

	return nil
}

func (f *feature) theResponseShouldMatchJson(message *messages.PickleStepArgument_PickleDocString) error {
	var expected, actual interface{}

	err := json.NewDecoder(strings.NewReader(message.Content)).Decode(&expected)
	if err != nil {
		return err
	}

	err = json.NewDecoder(f.resp.Body).Decode(&actual)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", expected, actual)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	f := &feature{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		f.resp = nil
	})

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, f.iSendRequestTo)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with json body:$`, f.iSendRequestToWithJsonBody)
	ctx.Step(`^the response code should be (\d+)$`, f.theResponseCodeShouldBe)
	ctx.Step(`^the response should match json:$`, f.theResponseShouldMatchJson)
}
