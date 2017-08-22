package controllers

import (
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

// Auth initiates the oauth2 authorization request to github
func (c App) Auth() revel.Result {
	user := c.currentUser()
	if user == nil {
		user = models.NewUser()
		c.Session["uid"] = fmt.Sprintf("%d", user.ID)
	}

	user.AuthEnv = app.Credentials.NewAuthEnvironment()
	authURL := user.AuthEnv.AuthCodeURL()
	revel.INFO.Printf("authURL= %s", authURL)
	return c.Redirect(authURL)
}

// AuthCallback handles the oauth2 authorization response and sets up a user
func (c App) AuthCallback() revel.Result {
	user := c.currentUser()
	if user == nil {
		revel.ERROR.Println("Invalid OAuth Callback")
		return c.Redirect("/")
	}

	state := c.Params.Query.Get("state")
	revel.INFO.Printf("DEBUG: AuthCallback state: %s", state)
	userState := user.AuthEnv.StateString
	if state != userState {
		revel.ERROR.Printf("Invalid OAuth State, expected '%s', got '%s'\n", userState, state)
		return c.Redirect("/")
	}

	code := c.Params.Query.Get("code")
	revel.INFO.Printf("DEBUG: AuthCallback code: %s", code)
	err := user.SetupUserAuth(code)
	if err != nil {
		revel.ERROR.Printf("Could not setup user authentication: %v", err)
		return c.Redirect("/")
	}
	revel.INFO.Printf("Successfully authenticated Github user: %s\n", user.GithubUsername)
	return c.Redirect("/workload")
}

func (c App) Workload() revel.Result {
	user := c.currentUser()
	if user == nil {
		revel.ERROR.Printf("User not setup correctly")
		return c.Redirect("/")
	}

	return c.Render(user)
}

func (c App) currentUser() *models.User {
	_, exists := c.Session["uid"]
	if !exists {
		return nil
	}

	var user *models.User
	uid, _ := strconv.ParseInt(c.Session["uid"], 10, 0)
	user = models.GetUser(int(uid))
	c.ViewArgs["user"] = user
	return user
}
