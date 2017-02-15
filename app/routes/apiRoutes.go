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
	router.GetFunc("/users/invites", api.AuthGetInvitees)

	// access tokens
	router.PostFunc("/access-tokens", api.AccessTokenCreate)

	// applications
	router.GetFunc("/applications", api.ApplicationIndex)
	router.GetFunc("/applications/:uuid", api.ApplicationShow)
	router.PostFunc("/applications", api.ApplicationCreate)
	router.PatchFunc("/applications/:uuid", api.ApplicationUpdate)

	// environments
	router.GetFunc("/environments", api.EnvironmentIndex)
	router.GetFunc("/environments/:uuid", api.EnvironmentShow)
	router.PatchFunc("/environments/:uuid", api.EnvironmentEdit)
	router.PostFunc("/environments", api.EnvironmentCreate)

	// versions
	router.GetFunc("/environments/:environment-uuid/versions", api.VersionIndex)
	router.GetFunc("/environments/:environment-uuid/versions/:version-uuid", api.VersionShow)
	router.PostFunc("/environments/:environment-uuid/versions", api.VersionCreate)
	router.PatchFunc("/environments/:environment-uuid/versions/:version-uuid", api.VersionUpdate)

	// environment config
	router.PostFunc("/environments/:environment-uuid/configuration", api.EnvironmentConfigCreate)

	return router
}

// apiRoutes registers the routes prefixed with /api/v1 where authentication is NOT required
func apiRoutes() *bone.Mux {
	router := bone.New()

	// JWT validation
	router.PostFunc("/auth/validate", api.ValidateToken)

	return router
}
