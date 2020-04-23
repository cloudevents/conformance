Feature: Kafka Protocol Binding

    The Kafka Protocol Binding for CloudEvents defines how events are mapped to Kafka messages.

    Background:
    Given Kafka Protocol Binding is supported

    Scenario: Binary content mode
    Given a Kafka message with payload:
        """
        {"message": "Hello World!"}
        """
    And Kafka headers:
        | key               | value                 |
        |    ce_specversion | 1.0                   |
        |             ce_id | 1234-1234-1234        |
        |           ce_type | com.example.someevent |
        |         ce_source | /mycontext/subcontext |
        |           ce_time | 2018-04-05T03:56:24Z  |
        |      content-type | application/json      |
    When parsed as Kafka message
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

    Scenario Outline: Structured content mode (<contentType>)
    Given a Kafka message with payload:
        """
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
    And Kafka headers:
        | key               | value         |
        |      content-type | <contentType> |
    When parsed as Kafka message
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