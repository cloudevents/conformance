@http
Feature: HTTP protocol binding

    The HTTP Protocol Binding for CloudEvents defines how events are mapped to HTTP 1.1 request and response messages.

    Background:
    Given HTTP Protocol Binding is supported

    Scenario Outline: Binary content mode (<contentType>)
    Given an HTTP request
        """
        POST /someresource HTTP/1.1
        Host: example.com
        ce-specversion: 1.0
        ce-type: com.example.someevent
        ce-time: 2018-04-05T03:56:24Z
        ce-id: 1234-1234-1234
        ce-source: /mycontext/subcontext
        Content-Type: <contentType>
        Content-Length: 33

        {
            "message": "Hello World!"
        }
        """
    When parsed as HTTP request
    Then the attributes are:
        | key               | value                           |
        |                id | 1234-1234-1234                  |
        |       specversion | 1.0                             |
        |              type | com.example.someevent           |
        |            source | /mycontext/subcontext           |
        |              time | 2018-04-05T03:56:24Z            |
        |   datacontenttype | <contentType>                   |
    And the data is equal to the following JSON:
        """
        {"message": "Hello World!"}
        """
    Examples:
    | contentType                     |
    | application/json                |
    | application/json; charset=utf-8 |

    Scenario Outline: Structured content mode (<contentType>)
    Given an HTTP request
        """
        POST /someresource HTTP/1.1
        Host: example.com
        Content-Type: <contentType>
        Content-Length: 266

        {
            "specversion": "1.0",
            "type": "com.example.someevent",
            "time": "2018-04-05T03:56:24Z",
            "id": "1234-1234-1234",
            "source": "/mycontext/subcontext",
            "datacontenttype": "application/json",
            "data": {
                "message": "Hello World!"
            }
        }
        """
    When parsed as HTTP request
    Then the attributes are:
        | key               | value                 |
        |                id | 1234-1234-1234        |
        |       specversion | 1.0                   |
        |              type | com.example.someevent |
        |            source | /mycontext/subcontext |
        |              time | 2018-04-05T03:56:24Z  |
        |   datacontenttype | application/json      |
    And the data is equal to the following JSON:
        """
        {"message": "Hello World!"}
        """
    Examples:
    | contentType                                 |
    | application/cloudevents+json                |
    | application/cloudevents+json; charset=utf-8 |