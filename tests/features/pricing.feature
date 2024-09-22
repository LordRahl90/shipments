Feature: Pricing
  Scenario: Get Price For Same Country and Small Weight
    Given I have country "se" and destination country of "se" and weight of 10.0
    When I get price"
    Then I should see the price of 100.00

  Scenario: Get Price For Different European Country and Small Weight
    Given I have country "se" and destination country of "dk" and weight of 10.0
    When I get price"
    Then I should see the price of 150.00

  Scenario: Get Price For Different Continental Country and Small Weight
    Given I have country "se" and destination country of "us" and weight of 10.0
    When I get price"
    Then I should see the price of 250.00