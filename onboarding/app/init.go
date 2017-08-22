package app

import (
	"github.com/revel/revel"
	onboarding "github.com/samsung-cnct/technical-on-boarding/onboarding/github/issues"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// Configs for onboard app loaded from conf/app.conf. These are required at startup.
	Configs = make(map[string]string)

	// AppSetup contains settings for the github workload job
	Setup *onboarding.SetupScheme

	// AppCredentials contains gitub app credentials
	Credentials *onboarding.Credentials
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(LoadConfigs)
	revel.OnAppStart(SetupScheme)
	revel.OnAppStart(SetupCredentials)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// Onboard specific configuration names
const (
	OnboardClientIdName     string = "onboard.client.id"
	OnboardClientSecretName string = "onboard.client.secret"
	OnboardOrgName          string = "onboard.org"
	OnboardRepoName         string = "onboard.repo"
	OnboardTasksFileName    string = "onboard.tasks.file"
	OnboardUserName         string = "onboard.user"
)

func LoadConfigs() {
	Configs[OnboardClientIdName] = revel.Config.StringDefault(OnboardClientIdName, "")
	Configs[OnboardClientSecretName] = revel.Config.StringDefault(OnboardClientSecretName, "")
	Configs[OnboardOrgName] = revel.Config.StringDefault(OnboardOrgName, "")
	Configs[OnboardRepoName] = revel.Config.StringDefault(OnboardRepoName, "")
	Configs[OnboardTasksFileName] = revel.Config.StringDefault(OnboardTasksFileName, "")

	for env, value := range Configs {
		if len(value) == 0 {
			revel.ERROR.Fatalf("The '%s' property is required on startup. check the conf/app.conf", env)
		}
	}
	revel.INFO.Printf("Configs Loaded")
}

func SetupScheme() {
	configFilename := Configs[OnboardTasksFileName]
	setup, err := onboarding.NewSetupScheme(configFilename, &Configs)
	if err != nil {
		revel.ERROR.Fatalf("Cannat create an onboarding setup scheme: %v", err)
	}
	Setup = setup
	revel.INFO.Printf("Scheme Setup")
}

func SetupCredentials() {
	Credentials = &onboarding.Credentials{
		ClientID:     Setup.ClientID,
		ClientSecret: Setup.ClientSecret,
		Scopes:       []string{"user", "repo", "issues", "milestones"},
	}
	revel.INFO.Printf("Credentials Setup")
}
