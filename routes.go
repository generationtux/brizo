package main

import (
	"github.com/go-zoo/bone"
)

// configureRoutes will configure the HTTP routes and appropriate handlers
func configureRoutes() *bone.Mux {
	router := bone.New()

	// Javascript UI
	router.HandleFunc("/app", uiHandler)
	router.HandleFunc("/app/*", uiHandler)

	// Healthz endpoint
	router.HandleFunc("/healthz", healthzHandler)

	return router
}
