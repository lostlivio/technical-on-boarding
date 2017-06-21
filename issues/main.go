package main

import (
	"io/ioutil"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"

	"fmt"
	"net/http"
	"os"

	"time"

	githuboauth "golang.org/x/oauth2/github"
)

type (
	// TaskEntry represents individual tasks to be assigned
	TaskEntry struct {
		Title       string
		Assignee    AssigneeIndirect `yaml:"assignee"`
		Description string           `yaml:"description,omitempty"`
	}

	AssigneeIndirect struct {
		GithubUsername string `yaml:"github_username"`
	}

	// SetupScheme represents the whole workload to be scheduled.
	SetupScheme struct {
		ClientID           string                      `yaml:"clientId"`
		ClientSecret       string                      `yaml:"clientSecret"`
		GithubOrganization string                      `yaml:"githubOrganization"`
		GithubRepository   string                      `yaml:"githubRepository"`
		Tasks              []TaskEntry                 `yaml:"tasks"`
		TaskOwners         map[string]AssigneeIndirect `yaml:"task_owners"`
	}
)

func (task *TaskEntry) String() *string {
	return &task.Assignee.GithubUsername
}

func (setup *SetupScheme) load(filename string) *SetupScheme {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &setup)
	// log.Printf("Result: #%v", setup)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Loaded configuration from file: %s", filename)
		// log.Printf("Found: %#v\n", setup) // super noisy.
	}
	return setup
}

func authHandleMain() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		const content = `<html><body><h2>Continue with <a href="/login">GitHub</a></h2></body></html>`

		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(content))
	}
}

func authHandleLogin(oauthStateString string, oauthConf *oauth2.Config) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		http.Redirect(writer, req, url, http.StatusTemporaryRedirect)
	}
}

func authGitHubCallback(oauthStateString string, oauthConf *oauth2.Config, successfulLoginHandler func(*github.Client) error) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		state := req.FormValue("state")
		if state != oauthStateString {
			fmt.Printf("invalid oauth state, expect '%s', got '%s' \n", oauthStateString, state)
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		code := req.FormValue("code")
		token, err := oauthConf.Exchange(oauth2.NoContext, code)
		if err != nil {
			fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		client := github.NewClient(oauthClient)
		user, _, err := client.Users.Get(oauth2.NoContext, "")
		if err != nil {
			fmt.Printf("client.Users.Get() failed with '%s'\n", err)
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Printf("Logged in as GitHub User: %s\n", *user.Login)
		integrationErr := successfulLoginHandler(client)
		if integrationErr != nil {
			fmt.Printf("Integration handler failed: %v\n", integrationErr)
		}

		http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
	}
}

func setupGitHubAuthComponents(setup *SetupScheme) (string, *oauth2.Config) {
	var (
		oauthConf = &oauth2.Config{
			ClientID:     setup.ClientID,
			ClientSecret: setup.ClientSecret,
			Scopes:       []string{"user", "repo", "issues", "milestones"},
			Endpoint:     githuboauth.Endpoint,
		}

		oauthStateString = "randomstring TODO"
	)
	return oauthStateString, oauthConf
}

func nvl(values ...*string) *string {
	for _, v := range values {
		if v == nil {
			continue
		}
		if len(*v) > 0 {
			return v
		}
	}
	return nil
}

// Target 3 weeks (rounding up) for onboarding completion.
// New hires starting on Mondays will effectively get 4 weeks.
func getMilestoneDueTime() time.Time {
	now := time.Now()
	offset := (time.Friday - now.Weekday())
	return now.AddDate(0, 0, 21+int(offset))
}

func resolveGitHubUser(client *github.Client, username string) *github.User {
	user, _, err := client.Users.Get(oauth2.NoContext, username)
	if err != nil {
		panic(err)
	}
	return user
}

func findExistingMilestone(client *github.Client, setup *SetupScheme, organizationName string, title *string) (*github.Milestone, error) {
	options := github.MilestoneListOptions{
		Sort:      "due_date",
		Direction: "desc",
	}

	existing, _, err := client.Issues.ListMilestones(oauth2.NoContext, organizationName, setup.GithubRepository, &options)

	if err != nil {
		return nil, err
	}

	for _, ms := range existing {
		if ms.GetTitle() == *title {
			return ms, nil // this one matched.
		}
	}

	return nil, nil // none matched
}

func createMilestone(client *github.Client, setup *SetupScheme, username string, organizationName string) (*github.Milestone, error) {

	user := resolveGitHubUser(client, username) // this will fail noisy if the user is not valid.
	username = user.GetLogin()

	title := fmt.Sprintf("Welcome @%s!", username)
	description := fmt.Sprintf("Let's setup up @%s for success. Here's what we need to cover...", username)
	duedate := getMilestoneDueTime()
	newMilestone := github.Milestone{
		Title:       &title,
		Description: &description,
		DueOn:       &duedate,
	}

	existing, err := findExistingMilestone(client, setup, organizationName, &title)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	milestone, _, err := client.Issues.CreateMilestone(oauth2.NoContext, organizationName, setup.GithubRepository, &newMilestone)
	return milestone, err
}

func createTickets(client *github.Client, setup *SetupScheme, username string, organizationName string, milestoneNumber int) ([]*github.Issue, error) {
	var issueResult ([]*github.Issue)

	for _, issueTemplate := range setup.Tasks {

		issue := github.IssueRequest{
			Title:     &issueTemplate.Title,
			Body:      &issueTemplate.Description,
			Assignee:  &issueTemplate.Assignee.GithubUsername,
			Milestone: &milestoneNumber,
		}

		result, _, err := client.Issues.Create(oauth2.NoContext, organizationName, setup.GithubRepository, &issue)
		if err != nil {
			log.Printf("GitHub Issues error: %v", err)
			log.Printf("Failed to create issue: %v [milestone: %d]", *issue.Title, *issue.Milestone)
			break
		} else {
			issueResult = append(issueResult, result)
		}
	}

	return issueResult, nil
}

func createProject(client *github.Client, setup *SetupScheme, username string, organizationName string) (*github.Project, error) {

	primaryName := "Samsung CNCT Onboarding"

	projects, _, err := client.Repositories.ListProjects(oauth2.NoContext, organizationName, setup.GithubRepository, nil)

	if err != nil {
		return nil, err
	}

	for _, aProject := range projects {
		if aProject.GetName() == primaryName {
			return aProject, nil
		}
	}

	options := github.ProjectOptions{
		Name: primaryName,
		Body: "Let's make it easy to for our newly joined team members to pick up speed.",
	}

	project, _, err := client.Repositories.CreateProject(oauth2.NoContext, organizationName, setup.GithubRepository, &options)

	if err != nil {
		return nil, err
	}

	columnNames := []string{"Backlog", "In Progress", "Blocked", "Done"}
	columnIds := []int{0, 0, 0, 0}

	for index, name := range columnNames {
		options := github.ProjectColumnOptions{
			Name: name,
		}
		column, _, err := client.Projects.CreateProjectColumn(oauth2.NoContext, project.GetID(), &options)

		if err != nil {
			return nil, err
		}

		columnIds[index] = column.GetID()
	}

	return project, err
}

func attachProjectIssue(client *github.Client, project *github.Project, columnName string, issue *github.Issue) (*github.ProjectCard, error) {
	cardOpts := github.ProjectCardOptions{
		ContentID:   issue.GetID(),
		ContentType: "issue",
	}

	// TODO: split column selection to a separate function
	// Should be followed by creating columns if target not found.

	columns, _, err := client.Projects.ListProjectColumns(oauth2.NoContext, project.GetID(), nil)

	var targetColumn *github.ProjectColumn
	targetColumn = nil

	for _, col := range columns {
		if col.GetName() == columnName {
			targetColumn = col
		}
	}

	log.Printf("Adding card for Issue #%d to column %d [%s]", issue.GetNumber(), targetColumn.GetID(), targetColumn.GetName())
	card, _, err := client.Projects.CreateProjectCard(oauth2.NoContext, targetColumn.GetID(), &cardOpts)

	if err != nil {
		return nil, err
	}

	// Maybe they have to be "moved" into the backlog column too?
	moveOpts := github.ProjectCardMoveOptions{
		Position: "bottom",
		ColumnID: targetColumn.GetID(),
	}

	_, err = client.Projects.MoveProjectCard(oauth2.NoContext, card.GetID(), &moveOpts)

	// TODO: nope, that didn't work. GitHub's UI still shows these issues as unaffiliated to the target project.
	// Issues will have to be manually added to the

	return card, err
}

func loginSuccessful(setup *SetupScheme, username string) func(*github.Client) error {
	return func(client *github.Client) error {

		log.Printf("Fetching repository %s/%s", setup.GithubOrganization, setup.GithubRepository)

		repo, _, err := client.Repositories.Get(oauth2.NoContext, setup.GithubOrganization, setup.GithubRepository)

		organizationName := repo.Owner.GetLogin()
		ptrOrganizatioName := &organizationName

		repoName := repo.GetFullName()
		ptrRepoName := &repoName

		if len(*ptrOrganizatioName) < 1 {
			ptrOrganizatioName = &setup.GithubOrganization
		}

		if len(*ptrRepoName) < 1 {
			ptrRepoName = &setup.GithubOrganization
		}

		if err != nil {
			return err
		}

		log.Printf("Setting up on-boarding workload on repo %s/%s", *ptrOrganizatioName, *ptrRepoName)

		// Prepare a milestone for the user's onboarding completion target.
		milestone, err := createMilestone(client, setup, username, *ptrOrganizatioName)

		if err != nil {
			return err
		}

		log.Printf("Created Milestone: %s <%s>", milestone.GetTitle(), milestone.GetURL())

		project, err := createProject(client, setup, username, *ptrOrganizatioName)

		if err != nil {
			return err
		}

		issues, err := createTickets(client, setup, username, *ptrOrganizatioName, milestone.GetNumber())

		if err != nil {
			return err
		}

		for _, i := range issues {
			log.Printf("Created Issue: #%d, '%s' [project: %s]", i.GetNumber(), i.GetTitle(), project.GetName())
			attachProjectIssue(client, project, "Backlog", i)
		}

		return err
	}
}

// Prepare GitHub Authetication Flow
func prepareGitHubAuthenticationFlow(setup *SetupScheme, username string, onSuccess func(*SetupScheme, string) func(*github.Client) error) {
	oauthStateString, oauthConf := setupGitHubAuthComponents(setup)
	log.Printf("Preparing GitHub integration: %v\n", setup)

	http.HandleFunc("/", authHandleMain())
	http.HandleFunc("/login", authHandleLogin(oauthStateString, oauthConf))
	http.HandleFunc("/github_oauth_cb", authGitHubCallback(oauthStateString, oauthConf, onSuccess(setup, username)))
	fmt.Print("Starting HTTP Service on http://127.0.0.1:7000\n")
	fmt.Println(http.ListenAndServe(":7000", nil))
}

func main() {
	username := os.Args[1]
	filename := os.Args[2]
	var setup SetupScheme

	prepareGitHubAuthenticationFlow(setup.load(filename), username, loginSuccessful)
}
