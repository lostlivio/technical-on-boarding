package models

import (
	"log"
	"math/rand"

	"github.com/google/go-github/github"
	onboarding "github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues"
	"golang.org/x/oauth2"
)

// User model object to manage user authentication
type User struct {
	ID             int
	GithubUsername string
	AccessToken    *oauth2.Token
	AuthEnv        *onboarding.AuthEnvironment
}

// SetupUserAuth sets the access token and retrieves the Github username
func (user *User) SetupUserAuth(code string) error {
	auth := user.AuthEnv
	token, err := auth.Config.Exchange(auth.Context, code)
	if err != nil {
		log.Printf("OAuth Exchange failed: %v", err)
		return err
	}
	user.AccessToken = token

	oauthClient := auth.Config.Client(auth.Context, token)
	githubClient := github.NewClient(oauthClient)
	githubUser, _, err := githubClient.Users.Get(auth.Context, "")
	if err != nil {
		log.Printf("Failed to get github user: %v", err)
		return err
	}
	user.GithubUsername = *githubUser.Login
	return nil
}

var db = make(map[int]*User)

// GetUser returns a User by id
func GetUser(id int) *User {
	return db[id]
}

// NewUser creates a new user
func NewUser() *User {

	user := &User{ID: rand.Intn(10000)}
	db[user.ID] = user
	return user
}
