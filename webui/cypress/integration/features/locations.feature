@core
Feature: Location Management
  # Scenarios on this feature need to be run in the defined order.

  Scenario: Create Location
    Given I logged in as admin
    When I go to locations page
    Then I should see page title "Locations"
    And I should see some entries
    When I click "Create" button
    Then I should see page title "New Location"
    When I fill in "name" with "test loc 1"
    And I fill in "warehouse" with "WH1"
    And I click "Submit" button
    And wait for "@locationCreate"
    Then I should see Location named "test loc 1"
    Then I should see Location in warehouse "WH1"
    When I click "Back" button
    Then I should see page title "Locations"
    And I should see 1 more entries

  Scenario: Cancel Create
    When I go to locations page
    Then I should see page title "Locations"
    And I should see some entries
    When I click "Create" button
    Then I should see page title "New Location"
    When I click "Cancel" button
    Then I should see page title "Locations"
    Then I should see 0 more entries

  # This scenario depends on the "Create Location" scenario.
  Scenario: Edit Location
    When I go to locations page
    And I click on link "test loc 1"
    And wait for location to load
    And I click "Edit" button
    Then I should see page title "Edit Location"
    And I fill in "warehouse" with "WH2"
    And I click "Submit" button
    And wait for "@locationUpdate"
    Then I should see Location named "test loc 1"
    Then I should see Location in warehouse "WH2"

  # This scenario depends on the "Create Location" scenario.
  Scenario: Delete Location
    When I go to locations page
    And I click on link "test loc 1"
    And wait for location to load
    And I click "Delete" button
    And wait for "@locationDelete"
    And I go to locations page
    Then I should see 1 fewer entries
