Feature: Inventory Transactions

  Background: Item and Location is ready
    Given I logged in
    And There is an item named "inventory test item"
    And There is a location named "inventory test location"

  Scenario Outline: Add transaction
    When I go to items page
    And I click on link "inventory test item"
    And wait for item to load
    Then I should see Item named "inventory test item"
    When I click on the plus icon button and wait
    And I select "inventory test item" in selector "item_id"
    And I select "inventory test location" in selector "location_id"
    And I check radio button <action>
    And I fill in "count" with <count>
    And I fill in "note" with <action>
    And I submit the inventory transaction
    Then I should see the latest transaction is for item "inventory test item" in location "inventory test location" for <diff>

    Examples:
    | action    | count | diff  |
    | "ADD"     | "10"  | "+10" |
    | "REMOVE"  | "5"   | "-5"  |
    | "RECOUNT" | "0"   | "[0]" |
