Feature: Create Shipment
  In order to create shipment, I need the base URL to connect with.

  Scenario: Create Small shipments within same country
    Given I am "Adewale James" With email "adewale@me.com"
    When I create a shipment with country "se" and destination country of "se" and weight of 10.0
    Then I should see the shipment with price of 100.00
    Then I should see the shipment with a non empty reference
