package routes

import (
	"github.com/generationtux/brizo/app/handlers/api"
	"github.com/go-zoo/bone"
)

// apiRoutes registers the routes prefixed with /api
func apiRoutes() *bone.Mux {
	router := bone.New()

	router.GetFunc("/login", api.AuthMainHandler)
	router.GetFunc("/o/auth/login/github", api.AuthGithubHandler)
	router.GetFunc("/o/auth/callback/github", api.AuthGithubCallbackHandler)

	return router
}
