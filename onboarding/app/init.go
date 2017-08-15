package app

import (
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	// AppConfigs for onboard app loaded from conf/app.conf. These are required at startup.
	AppConfigs map[string]string
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
	AppConfigs := map[string]string{
		OnboardClientIdName:     revel.Config.StringDefault(OnboardClientIdName, ""),
		OnboardClientSecretName: revel.Config.StringDefault(OnboardClientSecretName, ""),
		OnboardOrgName:          revel.Config.StringDefault(OnboardOrgName, ""),
		OnboardRepoName:         revel.Config.StringDefault(OnboardRepoName, ""),
		OnboardTasksFileName:    revel.Config.StringDefault(OnboardTasksFileName, ""),
	}

	for env, value := range AppConfigs {
		if len(value) == 0 {
			revel.ERROR.Fatalf("The '%s' property is required on startup. check the conf/app.conf", env)
		}
	}
}
