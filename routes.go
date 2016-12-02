package main

import (
	"github.com/go-zoo/bone"
)

// configureRoutes will configure the HTTP routes and appropriate handlers
func configureRoutes() *bone.Mux {
	router := bone.New()

	// Javascript UI
	router.HandleFunc("/ui", uiHandler)

	// Healthz endpoint
	router.HandleFunc("/healthz", healthzHandler)

	return router
}
