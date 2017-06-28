/*	This package abstracts authentication with GitHub's API.
	Primarily this extract OAuth integration from the rest of the business logic,
	enabling better testing of business logic without dependency on GitHub's actual
	service.
*/

package onboarding

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/satori/go.uuid"
)

type (

	// AuthenticationEnviron provides a simple model for OAuth context abstraction.
	AuthenticationEnviron struct {
		Context     context.Context
		Config      *oauth2.Config
		StateString string
	}

	// Credentials for a GitHub application, integrating with GitHub API.
	Credentials struct {
		ClientID     string
		ClientSecret string
		Scopes       []string
	}
)

// Login to GitHub and execute the given function on success.
func (creds *Credentials) Login(handler func(*github.Client, *context.Context) error) error {
	auth := NewAuthenticationEnvironment(creds.ClientID, creds.ClientSecret, creds.Scopes)
	return auth.authenticate(func(client *github.Client) error {
		return handler(client, &auth.Context)
	})
}

// NewAuthenticationEnvironment prepares an OAuth authenticated GitHub login environment.
func NewAuthenticationEnvironment(ClientID string, ClientSecret string, Scopes []string) *AuthenticationEnviron {
	var (
		authContext = oauth2.NoContext
		config      = oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			Scopes:       Scopes,
			Endpoint:     githuboauth.Endpoint,
		}
		oauthStateString = uuid.NewV4().String()
	)
	return &AuthenticationEnviron{
		Context:     authContext,
		Config:      &config,
		StateString: oauthStateString,
	}
}

func (auth *AuthenticationEnviron) prepareGitHubClient(tokenCode string) (*github.Client, error) {
	token, err := auth.Config.Exchange(auth.Context, tokenCode)
	if err != nil {
		log.Printf("OAuth Exchange failed: %v", err)
		return nil, err
	}
	oauthClient := auth.Config.Client(auth.Context, token)
	githubClient := github.NewClient(oauthClient)
	return githubClient, nil
}

// http handler "/"
func (auth *AuthenticationEnviron) handlerMain() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		const content = `<html><body><h2>Continue with <a href="/login">GitHub</a></h2></body></html>`

		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(content))
	}
}

// http handler "/login"
func (auth *AuthenticationEnviron) handlerLogin() func(http.ResponseWriter, *http.Request) {
	oauthStateString := auth.StateString
	oauthConf := *auth.Config
	return func(writer http.ResponseWriter, req *http.Request) {
		url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
		http.Redirect(writer, req, url, http.StatusTemporaryRedirect)
	}
}

// http handler "/github_oauth_cb"
func (auth *AuthenticationEnviron) handlerAuthenticated(onSuccess func(*github.Client) error) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		state := req.FormValue("state")
		if state != auth.StateString {
			log.Printf("Invalid OAuth State, expected '%s', got '%s'\n", auth.StateString, state)
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		code := req.FormValue("code")
		githubClient, err := auth.prepareGitHubClient(code)

		if err != nil {
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		user, _, err := githubClient.Users.Get(auth.Context, "")
		if err != nil {
			fmt.Printf("client.Users.Get() failed with '%s'\n", err)
			http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)
			return
		}

		// Successfully logged in.
		fmt.Printf("Logged in as GitHub User: %s\n", *user.Login)

		// Process the workload *after* redirecting the browser.
		defer onSuccess(githubClient)

		// Prevent timing out in the browser. Redirect immediately.
		http.Redirect(writer, req, "/", http.StatusTemporaryRedirect)

	}
}

// Sets up the HTTP server environment to negotiate OAuth via a browser, on behalf of the user.
func (auth *AuthenticationEnviron) authenticate(onSuccess func(*github.Client) error) error {

	http.HandleFunc("/", auth.handlerMain())
	http.HandleFunc("/login", auth.handlerLogin())
	http.HandleFunc("/github_oauth_cb", auth.handlerAuthenticated(onSuccess))

	log.Println("Please open your browser and navigate to http://127.0.0.1:7000 to authenticate...")
	log.Println(http.ListenAndServe(":7000", nil))

	return nil
}
