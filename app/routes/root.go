package routes

import (
	"net/http"

	"github.com/generationtux/brizo/app/handlers/web"
	"github.com/go-zoo/bone"
)

// BuildRouter configures the application router
func BuildRouter() *bone.Mux {
	r := mainRoutes()
	r.SubRoute("/api", apiRoutes())

	return r
}

// mainRoutes registers the routes for the user interface
func mainRoutes() *bone.Mux {
	router := bone.New()

	// Home page
	router.GetFunc("/", web.RootHandler)

	// Javascript UI
	router.GetFunc("/app", web.UIHandler)
	router.Get("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/dist/"))))

	// Healthz endpoint
	router.GetFunc("/healthz", web.HealthzHandler)

	return router
}
