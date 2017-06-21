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

1. [Register a New OAuth Application][2] in your GitHub profile.
    1. Specify the callback URL: `http://127.0.0.1:7000/github_oauth_cb`
    2. Copy the Client ID and Secret to a text file.
2. Build the Go program: `go build`
3. Execute the program: `./issues ${username} ./onboarding-issues.yaml`


## Known Issues

- [ ] Naive. Creates duplicate issues. Need to search existing by title & assignee. 
- [ ] Poorly tested. Refactoring in progress to support better testing.


## Testing

This workload relies heavily on the GitHub API, which also requires valid appliation tokens.

As of this time (initial commit) there is no testing-focused code. Refactoring is underway to improve this.
The minimum viable functionality was developed against a forked project in GitHub, effectively end-to-end testing only with real credentials.
(Pardon the mess. This is my first Go application too.)


[2]: https://github.com/settings/applications/new
[3]: https://github.com/settings/apps