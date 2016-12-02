package app

import (
	"github.com/go-zoo/bone"
)

// ConfigureRoutes will configure the HTTP routes and appropriate handlers
func ConfigureRoutes() *bone.Mux {
	router := bone.New()

	router.GetFunc("/", rootHandler)

	// Javascript UI
	router.GetFunc("/app", uiHandler)
	router.GetFunc("/app/*", uiHandler)

	// Healthz endpoint
	router.GetFunc("/healthz", healthzHandler)

	return router
}
