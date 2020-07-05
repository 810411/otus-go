# file: add_event.feature
Feature: add event
  In order to make sure that an event can added into calendar
  As an API user
  I need to add an event

  Scenario: successfully adding event
    When I send "POST" request to "http://localhost:8888/events" with json body:
      """
        {
          "title": "test",
          "datetime": "2020-07-01T20:07:01Z",
          "duration":1,
          "owner_id":1
        }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "event":{
          "id":1,
          "title": "test",
          "datetime": "2020-07-01T20:07:01Z",
          "duration":1,
          "description":"",
          "owner_id":1
        }
      }
      """

  Scenario: error - event's time busy
    When I send "POST" request to "http://localhost:8888/events" with json body:
      """
        {
          "title": "test",
          "datetime": "2020-07-01T20:07:01Z",
          "duration":1,
          "owner_id":1
        }
      """
    Then the response code should be 400
    And the response should match json:
      """
      {
        "status":"error",
        "error":{
          "message": "event's time busy"
        }
      }
      """

  Scenario: add event with same datetime but other user id
    When I send "POST" request to "http://localhost:8888/events" with json body:
      """
        {
          "title": "test",
          "datetime": "2020-07-01T20:07:01Z",
          "duration":1,
          "owner_id":2
        }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "event":{
          "id":3,
          "title": "test",
          "datetime": "2020-07-01T20:07:01Z",
          "duration":1,
          "description":"",
          "owner_id":2
        }
      }
      """

  Scenario: add event with new datetime
    When I send "POST" request to "http://localhost:8888/events" with json body:
      """
        {
          "title": "test",
          "datetime": "2020-07-08T20:07:01Z",
          "duration":1,
          "owner_id":1
        }
      """
    Then the response code should be 200
    And the response should match json:
      """
      {
        "status":"ok",
        "event":{
          "id":4,
          "title": "test",
          "datetime": "2020-07-08T20:07:01Z",
          "duration":1,
          "description":"",
          "owner_id":1
        }
      }
      """

  Scenario: wrong request method
    When I send "GET" request to "http://localhost:8888/events" with json body:
      """
        {
          "title": "test",
          "datetime": "2020-07-08T20:07:01Z",
          "duration":1,
          "owner_id":1
        }
      """
    Then the response code should be 405
