package routes

import (
	"github.com/generationtux/brizo/app/handlers/api"
	"github.com/go-zoo/bone"
)

// authAPIRoutes registers the routes prefixed with /api/v1 and protected by an authentication middleware
func authAPIRoutes() *bone.Mux {
	router := bone.New()

	// users
	router.PostFunc("/users", api.AuthCreateUser)

	// applications
	router.GetFunc("/applications", api.ApplicationIndex)
	router.GetFunc("/applications/:uuid", api.ApplicationShow)

	return router
}

// apiRoutes registers the routes prefixed with /api/v1 where authentication is NOT required
func apiRoutes() *bone.Mux {
	router := bone.New()

	// JWT validation
	router.PostFunc("/auth/validate", api.ValidateToken)

	return router
}
