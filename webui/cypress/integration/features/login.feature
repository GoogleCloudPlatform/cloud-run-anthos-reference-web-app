Feature: Avatar Photo
    The main page of the application will display an avatar photo when the user is
    logged in.

    Scenario: User is logged in
        When I log in
        Then my avatar image should be set