Feature: Inventory Transactions

  Background: Item and Location is ready
    Given I logged in as worker

  Scenario Outline: Add transaction
    When I go to items page
    And I click on test item
    And wait for item to load
    When I click on the plus icon button
    And I select test item
    And I select test location
    And I check radio button <action>
    And I fill in "count" with <count>
    And I fill in "note" with <action>
    And I submit the inventory transaction
    Then I should see the latest transaction is for test item in test location for <diff>

    Examples:
    | action    | count | diff  |
    | "ADD"     | "10"  | "+10" |
    | "REMOVE"  | "5"   | "-5"  |
    | "RECOUNT" | "0"   | "[0]" |
