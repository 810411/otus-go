# file: list_event.feature
Feature: list events
  In order to make sure that possible get events list from calendar
  As an API user
  I need to get list of event

  Scenario: list of day
    When I send "GET" request to "http://calendar:8888/events/day?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "events":[
          {"id":1,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":1},
          {"id":3,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":2}
        ]
      }
      """

  Scenario: list of week
    When I send "GET" request to "http://calendar:8888/events/week?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "events":[
          {"id":1,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":1},
          {"id":3,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":2}
        ]
      }
      """

  Scenario: list of month
    When I send "GET" request to "http://calendar:8888/events/month?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "events":[
          {"id":1,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":1},
          {"id":3,"title":"test","datetime":"2020-07-01T20:07:00Z","duration":1,"description":"","owner_id":2},
          {"id":4,"title":"test","datetime":"2020-07-08T20:07:00Z","duration":1,"description":"","owner_id":1}
        ]
      }
      """

  Scenario: list of eternity
    When I send "GET" request to "http://calendar:8888/events/eternity?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "status":"error",
        "error":{
          "message":"wrong period param (cases: day, week, month)"
        }
      }
      """

  Scenario: list of nothing
    When I send "GET" request to "http://calendar:8888/events?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 405

  Scenario: without datetime
    When I send "GET" request to "http://calendar:8888/events/day"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "status":"error",
        "error":{
          "message":"missing datetime query param"
        }
      }
      """

  Scenario: wrong request method
    When I send "PATCH" request to "http://calendar:8888/events/day?datetime=2020-07-01T00:00:00%2B0000"
    Then the response code should be 405

  Scenario: get something else
    When I send "GET" request to "http://calendar:8888/chicks"
    Then the response code should be 404

  Scenario: get from period where don't have events
    When I send "GET" request to "http://calendar:8888/events/day?datetime=2222-07-01T00:00:00%2B0000"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "events":null
      }
      """