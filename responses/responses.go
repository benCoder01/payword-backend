package responses

import (
	"net/http"

	"gitlab.com/benCoder01/payword-backend/db"

	"github.com/go-chi/render"
)

// UserResponse gibt öffentliche Daten eines Nutzers wieder
type UserResponse struct {
	Username string `json:"username"`
}

// NewUserResponse gibt den Benutzer an, der gefragt wurd.
func NewUserResponse(user *db.User) *UserResponse {
	return &UserResponse{
		Username: user.Username,
	}
}

func NewUserListResponse(users []db.User) []render.Renderer {
	list := []render.Renderer{}

	for _, user := range users {
		list = append(list, NewUserResponse(&user))
	}

	return list
}

// Render ...
func (ur *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (tr *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewTokenResponse(token string) *TokenResponse {
	return &TokenResponse{token}
}

type GameResponse struct {
	Name       string        `json:"name"`
	Admin      string        `json:"admin"`
	Categories []db.Category `json:"categories"`
	Members    []string      `json:"members"`
}

func (gr *GameResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewGameResponse(game *db.Game) *GameResponse {
	return &GameResponse{
		Name:       game.Name,
		Admin:      game.Admin,
		Categories: game.Categories,
		Members:    game.Members,
	}
}

func NewGameListResponse(games []db.Game) []render.Renderer {
	list := []render.Renderer{}

	for _, game := range games {
		list = append(list, NewGameResponse(&game))
	}

	return list
}

// Error Response Management

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

/*
Die unteren Methoden liefern ein Render Strukt, da hier auch die Methode Render implementiert.
*/

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

func ErrNotFound() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: 404,
		StatusText:     "Resource not found.",
	}
}
