package games

import (
	"errors"
	"net/http"

	"gitlab.com/benCoder01/payword-backend/db"

	"github.com/go-bongo/bongo"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"gitlab.com/benCoder01/payword-backend/requests"
	"gitlab.com/benCoder01/payword-backend/responses"
	"gitlab.com/benCoder01/payword-backend/token"
)

func Router(cors *cors.Cors) chi.Router {
	r := chi.NewRouter()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler)

	})

	// Protected
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler)

		r.Use(jwtauth.Verifier(token.TokenAuth))
		r.Use(jwtauth.Authenticator)

		// TODO: protected for Admin
		r.Post("/create", create)
		r.Get("/{name}", getGame)
		r.Get("/user/{name}", getGames)
		r.Post("/enter", addUser)
		r.Post("/leave", removeUser)

		r.Post("/increment", increment)
		r.Post("/decrement", decrement)

		r.Post("/category/add", addCategory)
		r.Post("/category/remove", removeCategory)
	})

	return r
}

func create(w http.ResponseWriter, r *http.Request) {
	gameReq := &requests.CreateGameRequest{}

	if err := render.Bind(r, gameReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	exists, err := db.GameExists(gameReq.Name)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if exists {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("Game name already exists")))
		return
	}

	game := db.Game{
		Name:       gameReq.Name,
		Admin:      gameReq.Admin,
		Categories: []db.Category{},
		Members:    []string{gameReq.Admin},
	}

	err = game.Save()

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(&game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}

}

func getGame(w http.ResponseWriter, r *http.Request) {
	gameName := chi.URLParam(r, "name")

	game, err := db.FindGameByName(gameName)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}

}

// getAllGames gibt alle Spiele, in denen der Nutzer eingetragen ist.
func getGames(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "name")

	gamesFromDb := db.FindAllGames()

	var games []db.Game

	// Spiele, bei denen nur der User mitspielt.
	for _, game := range gamesFromDb {
		for _, member := range game.Members {
			if member == username {
				games = append(games, game)
				break
			}
		}
	}

	if err := render.RenderList(w, r, responses.NewGameListResponse(games)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

// Ein Nutzer kann sich mit dem Namen des Spieles eintragen.
func addUser(w http.ResponseWriter, r *http.Request) {
	ctrlReq := &requests.EventControlRequest{}

	if err := render.Bind(r, ctrlReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(ctrlReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	// Search user in game
	if userInGame(game, ctrlReq.Username) {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("user already exists in game")))
		return
	}

	err = game.AddUser(ctrlReq.Username)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

// Testet, ob ein Benutzer bereits in einem Spiel enthalten ist.
func userInGame(game *db.Game, username string) bool {
	for _, member := range game.Members {
		if member == username {
			return true
		}
	}

	return false
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	ctrlReq := &requests.EventControlRequest{}

	if err := render.Bind(r, ctrlReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(ctrlReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	// Search user in game
	if !userInGame(game, ctrlReq.Username) {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("user not in game")))
		return
	}

	if game.Admin == ctrlReq.Username {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("the user is admin")))
		return
	}

	err = game.RemoveUser(ctrlReq.Username)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	categoryReq := &requests.CategoryControlRequest{}

	if err := render.Bind(r, categoryReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(categoryReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	if game.CategoryExists(categoryReq.Categoryname) {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("category alread exists")))
		return
	}

	err = game.AddCategory(categoryReq.Categoryname)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}

}

func removeCategory(w http.ResponseWriter, r *http.Request) {
	categoryReq := &requests.CategoryControlRequest{}

	if err := render.Bind(r, categoryReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(categoryReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	err = game.RemoveCategory(categoryReq.Categoryname)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func increment(w http.ResponseWriter, r *http.Request) {
	valueReq := &requests.ValueControlRequest{}

	if err := render.Bind(r, valueReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(valueReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	err = game.Increment(valueReq.Categoryname, valueReq.Username)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func decrement(w http.ResponseWriter, r *http.Request) {
	valueReq := &requests.ValueControlRequest{}

	if err := render.Bind(r, valueReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	game, err := db.FindGameByName(valueReq.Game)

	if err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	err = game.Decrement(valueReq.Categoryname, valueReq.Username)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewGameResponse(game)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}
