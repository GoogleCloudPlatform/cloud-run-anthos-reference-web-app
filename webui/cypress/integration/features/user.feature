@core
Feature: User Management

  Scenario: Manage User
    Given I logged in as admin
    When I go to users page
    Then I should see user with name "Test Admin" and role "admin"
    And I should see user with name "Test Worker" and role "worker"
    When I select role "admin" for user "Test Worker"
    Then I should see user with name "Test Worker" and role "admin"
    When I select role "worker" for user "Test Worker"
    Then I should see user with name "Test Worker" and role "worker"

  Scenario: You shall not be allowed
    Given I logged in as worker
    When I go to users page
    Then I should see user with name "Test Admin" and role "admin"
    And I should see user with name "Test Worker" and role "worker"
    And I should not be able to select roles
