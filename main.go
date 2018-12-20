package main

import (
	"fmt"
	"net/http"

	"gitlab.com/benCoder01/payword-backend/db"
	"gitlab.com/benCoder01/payword-backend/games"
	"gitlab.com/benCoder01/payword-backend/token"
	"gitlab.com/benCoder01/payword-backend/users"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	if err := db.Init(); err != nil {
		panic(err)
	}
	// configure Cors
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"LINK"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	})

	r.Use(
		middleware.RequestID, // Alle Requests bekommen eine ID
		middleware.Logger,
		cors.Handler,
		middleware.RealIP,
		middleware.Recoverer, // Fängt panics() ab.
		render.SetContentType(render.ContentTypeJSON),
	)

	// Generate Token Authenticator
	token.Init()

	// alle Requests an /user werden vom Router im user Package behandelt
	r.Mount("/users", users.Router(cors))
	r.Mount("/games", games.Router(cors))

	fmt.Println(http.ListenAndServe(":3333", r))
}
