Feature: Validate the input of a shipment

  Background: Validation rules
    Given "name" validation rules
    ```
    - Maximum: 30 characters
    - No numbers
    - Punctuation is allowed
    ```

    And "email" validation rules
    ```
    - Should be a valid email
    ```

    And "address" validation rules
    ```
    - Maximum: 100 characters
    ```

    And "country code" validation rules
    ```
    - Follow the ISO-3166-1 alpha-2 standard
    ```

    And "weight" validation rules
  ```
  - Weight is in kg
  - Maximum: 1000kg
  ```

  Scenario Outline: Create shipment with <sender_or_receiver> name: <name>
    Given a request to create a shipment with
      | <sender_or_receiver> - name | <name> |
    Then the returned error should have
      | message | <error> |

    Examples:
      | name                            | sender_or_receiver | error                                                                                                                                                | comment                            |
      |                                 | sender             |                                                                                                                                                      | empty name                         |
      |                                 | receiver           |                                                                                                                                                      | empty name                         |
      | User Example                    | sender             |                                                                                                                                                      | valid name                         |
      | User Example                    | receiver           |                                                                                                                                                      | valid name                         |
      | 1337 User                       | sender             | shipment was invalid: failed to validate sender: sender name is invalid: name is not valid, it contains numbers: 1337 User                           | numbers aren't allowed in the name |
      | 1337 User                       | receiver           | shipment was invalid: failed to validate receiver: receiver name is invalid: name is not valid, it contains numbers: 1337 User                       | numbers aren't allowed in the name |
      | AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA | sender             | shipment was invalid: failed to validate sender: sender name is invalid: name: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA is longer than the max length: 30     | too long name                      |
      | AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA | receiver           | shipment was invalid: failed to validate receiver: receiver name is invalid: name: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA is longer than the max length: 30 | too long name                      |

  Scenario Outline: Create shipment with <sender_or_receiver> email: <email>
    Given a request to create a shipment with
      | <sender_or_receiver> - email | <email> |
    Then the returned error should have
      | message | <error> |

    Examples:
      | email            | sender_or_receiver | error                                                                                                | comment              |
      |                  | sender             | shipment was invalid: failed to validate sender: sender email:  is not valid: invalid format         | empty email          |
      |                  | receiver           | shipment was invalid: failed to validate receiver: receiver email:  is not valid: invalid format     | empty email          |
      | user@example.com | sender             |                                                                                                      | valid email          |
      | user@example.com | receiver           |                                                                                                      | valid email          |
      | user             | sender             | shipment was invalid: failed to validate sender: sender email: user is not valid: invalid format     | email without domain |
      | user             | receiver           | shipment was invalid: failed to validate receiver: receiver email: user is not valid: invalid format | email without domain |

  Scenario Outline: Create shipment with <sender_or_receiver> address: <address>
    Given a request to create a shipment with
      | <sender_or_receiver> - address | <address> |
    Then the returned error should have
      | message | <error> |

    Examples:
      | address                                                                                               | sender_or_receiver | error                                                                                                                                                                                                     | comment          |
      |                                                                                                       | sender             |                                                                                                                                                                                                           | empty address    |
      |                                                                                                       | receiver           |                                                                                                                                                                                                           | empty address    |
      | Apt. Example 1A                                                                                       | sender             |                                                                                                                                                                                                           | valid address    |
      | Apt. Example 1A                                                                                       | receiver           |                                                                                                                                                                                                           | valid address    |
      | AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA | sender             | shipment was invalid: failed to validate sender: sender address: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA is longer than max length: 100     | too long address |
      | AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA | receiver           | shipment was invalid: failed to validate receiver: receiver address: AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA is longer than max length: 100 | too long address |

  Scenario Outline: Create shipment with <sender_or_receiver> country code: <country_code>
    Given a request to create a shipment with
      | <sender_or_receiver> - country code | <country_code> |
    Then the returned error should have
      | message | <error> |

    Examples:
      | country_code | sender_or_receiver | error                                                                                                                                                                                    | comment                |
      | S            | sender             | shipment was invalid: failed to validate sender: sender country code is invalid: could not find country by code: S, error: gountries error. Invalid code format: S                       | invalid                |
      | S            | receiver           | shipment was invalid: failed to validate receiver: receiver country code is invalid: could not find country by code: S, error: gountries error. Invalid code format: S                   | invalid                |
      | SE           | sender             |                                                                                                                                                                                          | valid country code     |
      | SE           | receiver           |                                                                                                                                                                                          | valid country code     |
      | SWE          | sender             | shipment was invalid: failed to validate sender: sender country code is invalid: country code: SWE is not of length: 2                                                                   | no support for alpha-3 |
      | SWE          | receiver           | shipment was invalid: failed to validate receiver: receiver country code is invalid: country code: SWE is not of length: 2                                                               | no support for alpha-3 |
      | ZZ           | sender             | shipment was invalid: failed to validate sender: sender country code is invalid: could not find country by code: ZZ, error: gountries error. Could not find country with code %s: ZZ     | invalid                |
      | ZZ           | receiver           | shipment was invalid: failed to validate receiver: receiver country code is invalid: could not find country by code: ZZ, error: gountries error. Could not find country with code %s: ZZ | invalid                |

  Scenario Outline: Create shipment with package weight: <weight>
    Given a request to create a shipment with
      | package - weight | <weight> |
    Then the returned error should have
      | message | <error> |

    Examples:
      | weight | error                                                                                               | comment          |
      | -1     | shipment was invalid: failed to validate package: package weight: -1 can't be below minimum: 0      | invalid          |
      | 0      |                                                                                                     | min valid weight |
      | 1000   |                                                                                                     | max valid weight |
      | 1001   | shipment was invalid: failed to validate package: package weight: 1001 can't be above maximum: 1000 | invalid          |