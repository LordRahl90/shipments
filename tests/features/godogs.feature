Feature: Eat some hotdogs
  In order to be happy and healthy,
  I need to be able to eat some hotdogs

  Scenario: Eat 5 of 12 hotdogs
    Given I have 12 hotdogs
    When I eat 5 hotdogs
    Then I should have 7 hotdogs left

  Scenario: Eat 10 of 12 hotdogs
    Given I have 12 hotdogs
    When I eat 10 hotdogs
    Then I should have 2 hotdogs left