package app

import (
	"net/http"

	"github.com/generationtux/brizo/auth"
	"github.com/go-zoo/bone"
)

// ConfigureRoutes will configure the HTTP routes and appropriate handlers
func ConfigureRoutes() *bone.Mux {
	router := bone.New()

	router.GetFunc("/", rootHandler)

	// Javascript UI
	router.GetFunc("/app", uiHandler)

	// Static Files
	router.Get("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("ui/dist/"))))

	// Healthz endpoint
	router.GetFunc("/healthz", healthzHandler)

	router.GetFunc("/login", auth.AuthMainHandler)
	router.GetFunc("/o/auth/login/github", auth.AuthGithubHandler)
	router.GetFunc("/o/auth/callback/github", auth.AuthGithubCallbackHandler)

	return router
}
