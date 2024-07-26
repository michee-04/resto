package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	jwtKey := os.Getenv("JWT_TOKEN_KEY")

	tokenAuth = jwtauth.New("HS256", []byte(jwtKey), nil)

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

	r.Route("/food", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create/{menuId}", controllers.CreateFodd)
		r.Get("/", controllers.GetFood)

		r.Route("/{foodId}", func(r chi.Router) {
			r.Get("/", controllers.GetfoodId)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateFood)
			r.With(provider.AdminOnly).Delete("/", controllers.Deletefood)
		})
	})

	r.Route("/invoice", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create/{userId}", controllers.CreateInvoice)
		r.Get("/", controllers.GetInvoice)

		r.Route("/{invoiceId}", func(r chi.Router) {
			r.Get("/", controllers.GetInvoiceId)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateInvoice)
			r.With(provider.AdminOnly).Delete("/", controllers.DeleteInvoice)
		})
	})

	r.Route("/note", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create", controllers.CreateNote)
		r.Get("/", controllers.GetNote)

		r.Route("/{noteId}", func(r chi.Router) {
			r.Get("/", controllers.GetNoteId)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateNote)
			r.With(provider.AdminOnly).Delete("/", controllers.DeleteNote)
		})
	})

	r.Route("/order", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create/{invoiceId}", controllers.CreateOrder)
		r.Get("/", controllers.GetOrder)

		r.Route("/{orderId}", func(r chi.Router) {
			r.Get("/", controllers.GetOrderId)
			r.With(provider.AdminOnly).Delete("/", controllers.DeleteOrder)
		})
	})

	r.Route("/order-item", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create", controllers.CreateOrderItem)
		r.Get("/", controllers.GetOrderItem)

		r.Route("/{orderItemId}", func(r chi.Router) {
			r.Get("/", controllers.GetOrderItemId)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateOrderItem)
			r.With(provider.AdminOnly).Delete("/", controllers.DeleteOrderIten)
		})
	})

	r.Route("/table", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.With(provider.AdminOnly).Post("/create", controllers.CreateTable)
		r.Get("/", controllers.GetTable)

		r.Route("/{tableId}", func(r chi.Router) {
			r.Get("/", controllers.GetTableId)
			r.With(provider.AdminOnly).Patch("/", controllers.UpdateTable)
			r.With(provider.AdminOnly).Delete("/", controllers.DeletTable)
		})
	})

	fmt.Printf("le serveur fonctionne sur http://localhost%s", port)

	http.ListenAndServe(port, r)
}
