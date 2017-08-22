/*	This package abstracts authentication with GitHub's API.
	Primarily this extract OAuth integration from the rest of the business logic,
	enabling better testing of business logic without dependency on GitHub's actual
	service.
*/

package onboarding

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type (
	// AuthEnvironment provides a simple model for OAuth context abstraction.
	AuthEnvironment struct {
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

// NewAuthEnvironment prepares a new OAuth2 authenticated GitHub login environment.
func (creds *Credentials) NewAuthEnvironment() *AuthEnvironment {
	var (
		authContext = oauth2.NoContext
		config      = oauth2.Config{
			ClientID:     creds.ClientID,
			ClientSecret: creds.ClientSecret,
			Scopes:       creds.Scopes,
			Endpoint:     githuboauth.Endpoint,
		}
		oauthStateString = uuid.NewV4().String()
	)
	return &AuthEnvironment{
		Context:     authContext,
		Config:      &config,
		StateString: oauthStateString,
	}
}

func (auth *AuthEnvironment) AuthCodeURL() string {
	oauthStateString := auth.StateString
	oauthConf := *auth.Config
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	return url
}
