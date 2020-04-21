Feature: Item Management

  Scenario: Create Item
    Given I logged in
    When I go to Items page
    Then I should see some items as start
    When I click "Create" button
    And I fill in "name" with "test item 1"
    And I fill in "description" with "description of test item"
    And I click Submit button
    Then I should see Item with "test item 1" and "description of test item"
    When I click "Back" button
    Then I should see 1 more items than start

  Scenario: Edit Item
    Given I logged in
    When I go to Items page
    And I click on link "test item 1"
    And I click "Edit" button
    And I fill in "description" with "edited description of test item"
    And I click Submit button
    Then I should see Item with "test item 1" and "edited description of test item"

  Scenario: Delete Item
    Given I logged in
    When I go to Items page
    And I click on link "test item 1"
    And I click "Delete" button
    And I go to Items page
    Then I should see 0 more items than start
