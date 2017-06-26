# Onboarding Workflow Template




## Goals

- Automate the preparation of a new SDSA employee's onboarding, as a list of tasks, represented as issues in GitHub.
- Provide a functional example of a Go program, integrating with the GitHub API.

### Functional Requirements

- Loads a template of "tasks" to be assigned to a new-hire. 
- Creates a Milestone and Project in GitHub. 
- Creates Issues in GitHub to represent tasks, and links them to Milestone and Project.
- Assigns those Issues to the new-hire.

## Usage

To run this application, you need credentials for the `SDSA onboarding workflow` application. 
Please ping `@here` in the Slack channel `#team-tooltime` if you need these credentials.

```shell
sh setup.sh  # installs dependencies, etc
make test    # runs Golang unit tests/etc
make build   # prepares an executable
./prepare_workload ./onboarding-issues.yaml  # executes the configured workload against GitHub
```



## Development & Testing

This workload relies heavily on the GitHub API, which also requires valid appliation tokens.

To facilitate testing this projects includes a fairly robust mock of the GitHub API client, and relies on
interfaces and proxy methods in several other points to allow the business logic to operate against a local
testing environment without reaching GitHub's API service.

Per golang's convention, tests are found in files ending with `_test.go`.





[2]: https://github.com/settings/applications/new
[3]: https://github.com/settings/apps