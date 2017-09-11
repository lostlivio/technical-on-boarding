package main

import (
	"github.com/samsung-cnct/technical-on-boarding/onboarding/app"
	_ "github.com/samsung-cnct/technical-on-boarding/onboarding/app/jobs/github"
)

func main() {
	// manually init app
	app.LoadConfigs()
	app.SetupScheme()
	app.SetupCredentials()

	// start job
	// events := make(<-chan Event)
	// job := github.GenerateProject{id,...,...,events}
	// go job.Run()
	// for channel switch on events
	//   fmt.Printf(event.Text)
}
