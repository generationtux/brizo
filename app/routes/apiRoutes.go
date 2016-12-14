package routes

import (
	"github.com/generationtux/brizo/app/handlers/api"
	"github.com/go-zoo/bone"
)

// apiRoutes registers the routes prefixed with /api
func apiRoutes() *bone.Mux {
	router := bone.New()

	// users
	router.PostFunc("/users", api.AuthCreateUser)

	// applications
	router.GetFunc("/applications", api.Applications)
	router.GetFunc("/applications/:name", api.Application)

	return router
}
