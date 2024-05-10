Feature: Update order production state
  In order to update the state of an order
  As a user working at the kitchen
  I want to be able to update the state of an order

  Scenario: Update order state
    Given I have an order
    When I update the order state to "Processing"
    Then the order state should be "Processing"
