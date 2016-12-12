package app

import (
	"net/http"

	"github.com/generationtux/brizo/app/handlers/api"
	"github.com/generationtux/brizo/app/handlers/web"
	"github.com/go-zoo/bone"
)

// ConfigureRoutes will configure the HTTP routes and appropriate handlers
func ConfigureRoutes() *bone.Mux {
	router := bone.New()

	router.GetFunc("/", web.RootHandler)

	// Javascript UI
	router.GetFunc("/app", web.UiHandler)

	// Static Files
	router.Get("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("ui/dist/"))))

	// Healthz endpoint
	router.GetFunc("/healthz", web.HealthzHandler)

	router.GetFunc("/login", api.AuthMainHandler)
	router.GetFunc("/o/auth/login/github", api.AuthGithubHandler)
	router.GetFunc("/o/auth/callback/github", api.AuthGithubCallbackHandler)

	// JWT handlers
	router.GetFunc("/get-token", web.GetTokenHandler)

	return router
}
