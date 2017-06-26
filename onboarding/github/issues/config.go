package onboarding

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type (
	// TaskEntry represents individual tasks to be assigned
	TaskEntry struct {
		Title       string
		Assignee    indirectAssignee `yaml:"assignee"`
		Description string           `yaml:"description,omitempty"`
	}

	indirectAssignee struct {
		GithubUsername string `yaml:"github_username"`
	}

	// SetupScheme represents the whole workload to be scheduled.
	SetupScheme struct {
		ClientID           string                      `yaml:"clientId"`
		ClientSecret       string                      `yaml:"clientSecret"`
		GithubOrganization string                      `yaml:"githubOrganization"`
		GithubRepository   string                      `yaml:"githubRepository"`
		Tasks              []TaskEntry                 `yaml:"tasks"`
		TaskOwners         map[string]indirectAssignee `yaml:"task_owners"`
	}
)

func (assignee *indirectAssignee) String() string {
	return assignee.GithubUsername
}

func (setup *SetupScheme) ingest(data []byte) error {

	err := yaml.Unmarshal(data, &setup)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (setup *SetupScheme) load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return setup.ingest(data)
}
