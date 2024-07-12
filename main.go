package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/michee-04/resto/controllers"
)

const port = ":8080"

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.CreateUser)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/", controllers.GetUser)
	})

	fmt.Printf("le serveur fonctionne sur http://localhost%s", port)

	http.ListenAndServe(port, r)
}