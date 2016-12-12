package routes

import (
	"github.com/go-zoo/bone"
)

// apiRoutes registers the routes prefixed with /api
func apiRoutes() *bone.Mux {
	router := bone.New()

	// users
	router.PostFunc("/users", api.AuthCreateUser)

	return router
}
