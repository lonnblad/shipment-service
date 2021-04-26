Feature: Calculate the price of a shipment

  Background: Price rules
    Given "region" price rules
    ```
    - Nordic region (Sweden, Norway, Denmark, Finland), the price is multiplied by 1
    - EU region, the price is multiplied by 1.5
    - Outside the EU, the price is multiplied by 2.5
    ```

    And "weight-class" price rules
    ```
    - Small (0 - 10kg): 100sek
    - Medium (10 - 25kg): 300sek
    - Large (26 - 50kg): 500sek
    - Huge (51 - 1000kg): 2000sek
    ```

    And price equation "{region}*{weight_class}"

  Scenario Outline: Create shipment for package: <package (kg)>, sender: <sender>
    Given a request to create a shipment with
      | sender - country code | <sender>       |
      | package - weight      | <package (kg)> |
    Then the returned shipment should have
      | package - price | <price (SEK)> |

    Examples:
      | sender | package (kg) | price (SEK) |
      | US     | 45           | 1250        |
      | SE     | 45           | 500         |
