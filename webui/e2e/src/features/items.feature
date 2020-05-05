Feature: Item Management
  # Scenarios on this feature need to be run in the defined order.

  Scenario: Create Item
    Given I logged in
    When I go to "items" page
    Then I should see page title "Items"
    And I should see some entries
    When I click "Create" button
    Then I should see page title "New Item"
    When I fill in "name" with "test item 1"
    And I fill in "description" with "description of test item"
    And I click "Submit" button
    Then I should see Item named "test item 1"
    And I should see Item description as "description of test item"
    When I click "Back" button
    Then I should see page title "Items"
    And I should see 1 more entries

  Scenario: Cancel Create
    When I go to "items" page
    Then I should see page title "Items"
    And I should see some entries
    When I click "Create" button
    Then I should see page title "New Item"
    When I click "Cancel" button
    Then I should see page title "Items"
    Then I should see 0 more entries

  # This scenario depends on the "Create Item" scenario.
  Scenario: Edit Item
    When I go to "items" page
    And I click on link "test item 1"
    And I click "Edit" button
    Then I should see page title "Edit Item"
    When I fill in "description" with "edited description of test item"
    And I click "Submit" button
    Then I should see Item named "test item 1"
    And I should see Item description as "edited description of test item"

  # This scenario depends on the "Create Item" scenario.
  Scenario: Delete Item
    When I go to "items" page
    And I click on link "test item 1"
    And I click "Delete" button
    And I go to "items" page
    Then I should see 1 fewer entries
