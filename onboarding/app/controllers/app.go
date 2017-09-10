package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/revel/revel"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/jobs"
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app/models"
	"golang.org/x/net/websocket"
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

	// In order to select between websocket messages and subscription events, we
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

	// Emit fake events every 5 seconds
	fakeEvents := make(chan jobs.Event)
	go func() {
		var cnt int

		for {
			cnt++
			msg := fmt.Sprintf("fake event %d", cnt)
			fakeEvents <- jobs.NewEvent(user.ID, "progress", msg)
			time.Sleep(5 * time.Second)
		}
	}()

	// Now listen for new events from either the websocket or the job.
	for {
		select {
		case event := <-fakeEvents:
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
			revel.INFO.Printf("Sending event: %v", event)
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
