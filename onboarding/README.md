# Onboarding Workflow Template



## Goals

- Automate the preparation of a new CNCT employee's onboarding, as a list of tasks, represented as issues in GitHub.
- Provide a functional example of a Go program, integrating with the GitHub API.

### Functional Requirements

- Loads a template of "tasks" to be assigned to a new-hire. 
- Creates a Milestone and Project in GitHub. 
- Creates Issues in GitHub to represent tasks, and links them to Milestone and Project.
- Assigns those Issues to the new-hire.


## Usage

To run this application, you need credentials for the `SDSA onboarding workflow` application. 
Please ping `@here` in the Slack channel `#team-tooltime` if you need these credentials.

The recommended approach is to execute this through Docker, like so:

```shell
docker build -t "samsung-cnct-onboarding:local" ./
docker run --rm -it -p 7000:7000 \
    -e 'GITHUB_CLIENT_ID={clientid}' \
    -e 'GITHUB_CLIENT_SECRET={clientsecret}'\
    -e 'GITHUB_ORG=samsung-cnct' -e 'GITHUB_REPO=technical-on-boarding'\
    -e 'GITHUB_USER={your_github_username}'\
    "samsung-cnct-onboarding:local"
```

This will start a local HTTP server, at [127.0.0.1:7000](http://127.0.0.1:7000/). Open this URL
in your browser, and log into GitHub _as yourself_. Once authenticated, the application will
return you to the local web server, and set up the workload in a GitHub project. 

The results of this application are logged to your terminal. 
**TODO** Eventually we will render a status page to the browser with results of the application's workload.

#### Alternate listening IP addresses.

*NOTE*: In some Docker configurations this may listen on a different IP address.
Typically this is most easily found via:

```shell
env | grep DOCKER_HOST  # hopefully not empty...
```

Or potentially:

```shell
docker run --net=host codenvy/che-ip 
```

If either of the above indicate another IP address, try reaching that on port 7000.

### Running Locally (in your desktop)

If somehow you prefer to build and run this in Go on your desktop environment,
rather than a preconfigured Docker container, here's how...

Please ensure that [Go is properly set up](./SETTINGUPGO.md) first.


```shell
make setup   # installs dependencies, etc
make test    # runs Golang unit tests/etc
make build   # prepares an executable

# app credentials required
export GITHUB_CLIENT_ID="{clientid}" \ 
    GITHUB_CLIENT_SECRET="{clientsecret}" \
    GITHUB_REPO="technical-on-boarding" \
    GITHUB_ORG="samsung-cnct" \
    GITHUB_USER="{your_github_username}"

# execute the configured workload against GitHub
./prepare_workload ./onboarding-issues.yaml  
```

## Development & Testing

This workload relies heavily on the GitHub API, which also requires valid appliation tokens.

To facilitate testing this projects includes a fairly robust mock of the GitHub API client, and relies on
interfaces and proxy methods in several other points to allow the business logic to operate against a local
testing environment without reaching GitHub's API service.

Per golang's convention, tests are found in files ending with `_test.go`.





[2]: https://github.com/settings/applications/new
[3]: https://github.com/settings/apps