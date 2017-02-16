@datacenter @datacenter_create @datacenter_create_azure
Feature: Ernest datacenter create

  Scenario: Non logged azure datacenter creation
    Given I setup ernest with target "https://ernest.local"
    And I logout
    When I run ernest with "datacenter create azure"
    Then The output should contain "You should specify the datacenter name"
    When I run ernest with "datacenter create azure tmp_datacenter"
    Then The output should contain "You're not allowed to perform this action, please log in"

  Scenario: Logged user azure datacenter creation
    Given I setup ernest with target "https://ernest.local"
    And the datacenter "tmp_datacenter" does not exist
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "datacenter create azure"
    Then The output should contain "You should specify the datacenter name"
    When I run ernest with "datacenter create azure tmp_datacenter"
    Then The output should contain "Please, fix the error shown below to continue"
    And The output should contain "- Specify a valid subscription id key with --subscription_id flag"
    When I run ernest with "datacenter create azure --subscription_id subscription --client_id client --client_secret secret --tenant_id tenant tmp_datacenter"
    Then The output should contain "Datacenter 'tmp_datacenter' successfully created"
    When I run ernest with "datacenter list"
    Then The output should contain "tmp_datacenter"
    Then The output should contain "tmp_region"
    Then The output should contain "azure"

  Scenario: Adding an already existing azure datacenter
    Given I setup ernest with target "https://ernest.local"
    And the datacenter "tmp_datacenter" does not exist
    And I'm logged in as "usr" / "pwd"
    When I run ernest with "datacenter create azure --subscription_id subscription --client_id client --client_secret secret --tenant_id tenant tmp_datacenter"
    And I run ernest with "datacenter create azure --subscription_id subscription --client_id client --client_secret secret --tenant_id tenant tmp_datacenter"
    Then The output should contain "Datacenter 'tmp_datacenter' already exists, please specify a different name"


