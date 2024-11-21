package router

import (
	"net/http"
	"pomodoro-rpg-api/presentation/handler"
	"pomodoro-rpg-api/presentation/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type HandlerDependencies struct {
	AuthHandler    handler.AuthHandler
	AccountHandler handler.AccountHandler
	TimeHandler    handler.TimeHandler
}

func New(deps HandlerDependencies, authenticator *middleware.Authenticator) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})

	r.Get("/is-auth", deps.AuthHandler.IsAuth)
	r.Post("/signup", deps.AuthHandler.SignUp)
	r.Post("/signup/confirm", deps.AuthHandler.ConfirmSignUp)
	r.Post("/signin", deps.AuthHandler.SignIn)
	r.Post("/forgot-password", deps.AuthHandler.ForgotPassword)
	r.Post("/forgot-password/confirm", deps.AuthHandler.ConfirmForgotPassword)

	r.Group(func(r chi.Router) {
		r.Use(authenticator.Middleware)

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", deps.AccountHandler.Get)
			r.Put("/", deps.AccountHandler.Update)
		})

		r.Post("/signout", deps.AuthHandler.SignOut)
		r.Post("/times", deps.TimeHandler.Create)
		r.Post("/change-password", deps.AuthHandler.ChangePassword)
	})

	return r
}
