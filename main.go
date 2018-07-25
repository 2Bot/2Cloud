package main

import (
	"github.com/go-chi/chi"
)

func main() {
	// creates a chi router
	r := chi.NewRouter()

	r.Route("/{userID}", func(r chi.Router) {
		// gets data for {userID}'s instance
		r.Get("/", getData)
		// creates an instance for {userID}
		r.Post("/create", createContainer)
		// updates {userID}'s instance
		r.Post("/update", updateContainer)
		// updates the settings for {userID}'s instance. Settings sent in body
		r.Put("/settings", updateSettings)
	})
}

func getData() {}

func createContainer() {}

func updateContainer() {}

func updateSettings() {}
