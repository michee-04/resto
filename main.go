package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/michee-04/resto/controllers"
	"github.com/michee-04/resto/provider"
)

const port = ":8080"

var tokenAuth *jwtauth.JWTAuth

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tokenAuth = jwtauth.New("HS256", []byte("ksQD5adHXZ-5SSJCupcHwBzDi6q5kfr5hdU7Eq5tMmo"), nil)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.CreateUser)
		r.Get("/verify", controllers.VerifyHandler)
		r.Get("/google/login", controllers.GoogleLogin)
		r.Get("/google/callback", controllers.GoogleCallback)
		r.Post("/login", controllers.LoginHandler)
		r.Post("/forgot-password", controllers.ForgotPasswordHandler)
		r.Get("/reset-password-email", controllers.ResetPasswordEmail)
		r.Post("/reset-password", controllers.ResetPasswordHandler)
		r.Patch("/set-password", controllers.SetPasswordHandler)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/", controllers.GetUser)

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Get("/", controllers.GetUserById)
			r.Post("/", controllers.LogoutUser)
			r.Patch("/", controllers.UpddateUser)
			r.Delete("/", controllers.DeleteUser)
		})
	})

	r.Route("/menu", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create", controllers.CreateMenu)
		r.Get("/", controllers.GetMenu)
		
		r.Route("/{menuId}", func(r chi.Router) {
			r.Get("/", controllers.GetMenuById)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateMenu)
			r.With(provider.AdminOnly).Delete("/", controllers.Deletemenu)
		})
	})

	fmt.Printf("le serveur fonctionne sur http://localhost%s", port)

	http.ListenAndServe(port, r)
}
