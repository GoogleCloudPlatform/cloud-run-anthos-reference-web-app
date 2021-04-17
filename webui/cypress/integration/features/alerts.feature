@alerts
Feature: Alerts

  Background: Item and Location is ready
    Given I logged in as worker

  Scenario Outline: Transaction which results in alert
    When I go to items page
    And I click on test item
    And wait for item to load
    When I click on the plus icon button and wait
    And I select test item
    And I select test location
    And I check radio button <action>
    And I fill in "count" with <count>
    And I submit the inventory transaction
    And I wait 30 seconds
    And I go to alerts page
    Then I should see the latest alert is for test item contains <text>

   Examples:
    | action    | count  | text                                   |
    | "RECOUNT" | "0"    | "Low total inventory for item: 0."     |
    | "ADD"     | "2000" | "High total inventory for item: 2000." |

  # This scenario depends on the "Transaction which results in alert" scenario.
  Scenario: Delete Alert
    When I go to alerts page
    Then I should see some entries
    When I dismiss the latest alert
    Then I should see 1 fewer entries