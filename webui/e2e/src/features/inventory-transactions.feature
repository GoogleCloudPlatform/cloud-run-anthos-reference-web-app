Feature: Inventory Transactions

  Background: Item and Location is ready
    Given I logged in
    And There is an item named "inventory test item"
    And There is a location named "inventory test location"
    When I go to "items" page
    And I click on link "inventory test item"
    Then I should see Item named "inventory test item"

  Scenario: Add transaction
    When I click on the plus icon button
    And I select "inventory test item" in selector "item_id"
    And I select "inventory test location" in selector "location_id"
    And I check radio button "ADD"
    And I fill in "count" with "10"
    And I fill in "note" with "test adding"
    And I click "Submit" button
    Then I should see the latest transaction is for item "inventory test item" in location "inventory test location" for "+10"


  Scenario: Remove transaction
    When I click on the plus icon button
    And I select "inventory test item" in selector "item_id"
    And I select "inventory test location" in selector "location_id"
    And I check radio button "REMOVE"
    And I fill in "count" with "5"
    And I fill in "note" with "test removing"
    And I click "Submit" button
    Then I should see the latest transaction is for item "inventory test item" in location "inventory test location" for "-5"

  Scenario: Recount transaction
    When I click on the plus icon button
    And I select "inventory test item" in selector "item_id"
    And I select "inventory test location" in selector "location_id"
    And I check radio button "RECOUNT"
    And I fill in "count" with "0"
    And I fill in "note" with "test recount"
    And I click "Submit" button
    Then I should see the latest transaction is for item "inventory test item" in location "inventory test location" for "[0]"
