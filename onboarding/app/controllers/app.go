package controllers

import (
	"fmt"
	"strconv"

	"github.com/revel/revel"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/jobs"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/jobs/onboarding"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/models"
	"golang.org/x/net/websocket"
)

// App for revel controller
type App struct {
	*revel.Controller
}

// Version endpoint to retrieve and serve app version.
// Can be used for an application readiness check.
func (c App) Version() revel.Result {
	return c.RenderJSON(app.SemanticVersion)
}

// Index of web app
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

	auth := app.Credentials.NewAuthEnvironment()
	authURL := auth.AuthCodeURL()
	user.AuthEnv = auth
	return c.Redirect(authURL)
}

// AuthCallback handles the oauth2 authorization response and sets up a user
func (c App) AuthCallback() revel.Result {
	user := c.currentUser()
	if user == nil {
		revel.ERROR.Println("Invalid OAuth Callback")
		return c.Redirect("/")
	}

	auth := user.AuthEnv
	state := c.Params.Query.Get("state")
	userState := auth.StateString
	if state != userState {
		revel.ERROR.Printf("Invalid OAuth State, expected '%s', got '%s'\n", userState, state)
		return c.Redirect("/")
	}

	code := c.Params.Query.Get("code")
	_, err := auth.SetupAccessToken(code)
	if err != nil {
		revel.ERROR.Printf("Could not get access token for user: %v", err)
		return c.Redirect("/")
	}
	user.Username = auth.GithubUsername()

	revel.INFO.Printf("Successfully authenticated Github user: %s\n", user.Username)
	return c.Redirect("/workload")
}

// Workload handles the initial workload page rendering
func (c App) Workload() revel.Result {
	user := c.currentUser()
	if user == nil {
		revel.ERROR.Printf("User not setup correctly")
		return c.Redirect("/")
	}

	return c.Render(user)
}

// WorkloadSocket handles the websocket connection for workload events
func (c App) WorkloadSocket(ws *websocket.Conn) revel.Result {
	if ws == nil {
		revel.ERROR.Printf("Websocket not intialized")
		return nil
	}
	user := c.currentUser()
	if user == nil {
		revel.ERROR.Printf("User not setup correctly")
		return c.Redirect("/")
	}

	// In order to select between websocket messages and job events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// setup and execute job
	events := make(chan jobs.Event)
	job := onboarding.GenerateProject{
		ID:      user.ID,
		Setup:   app.Setup,
		AuthEnv: user.AuthEnv,
		New:     events,
	}
	jobs.StartJob(job)

	// Now listen for new events from either the websocket or the job.
	for {
		select {
		case event, ok := <-events:
			if !ok {
				// Completed job events
				revel.INFO.Printf("The job has completed")
				return nil
			}
			revel.INFO.Printf("Sending event: %v", event)
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				revel.INFO.Printf("The user '%s' has disconnected", user.Username)
				return nil
			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}
			revel.INFO.Printf("Recieved: " + msg)
		}
	}
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
