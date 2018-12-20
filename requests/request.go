package requests

import (
	"errors"
	"net/http"
)

// UserRequest bezeichnet eine Anfrage mit Benutzerdaten
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Bind ist eine Methode, die Durch das Interface aus Render bedingt ist.
// Sie dient als Middleware beim Parsen der Request.
func (ur *UserRequest) Bind(r *http.Request) error {
	if ur == nil {
		return errors.New("missing required user fields")
	}

	return nil

}

type PasswordChangeRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
}

func (pcr *PasswordChangeRequest) Bind(r *http.Request) error {
	if pcr == nil {
		return errors.New("missing required user fields")
	}

	return nil
}

type CreateGameRequest struct {
	Name  string `json:"name"`
	Admin string `json:"admin"`
}

func (cgr *CreateGameRequest) Bind(r *http.Request) error {
	if cgr == nil {
		return errors.New("missing required user fields")
	}

	return nil
}

// Request, um einem Spiel beizutreten, oder es zu verlassen.
type EventControlRequest struct {
	Username string `json:"username"`
	Game     string `json:"game"`
}

func (cgr *EventControlRequest) Bind(r *http.Request) error {
	if cgr == nil {
		return errors.New("missing required user fields")
	}

	return nil
}

type CategoryControlRequest struct {
	Game         string `json:"game"`
	Categoryname string `json:"categoryname"`
}

func (cgr *CategoryControlRequest) Bind(r *http.Request) error {
	if cgr == nil {
		return errors.New("missing required user fields")
	}

	return nil
}

type ValueControlRequest struct {
	Game         string `json:"game"`
	Categoryname string `json:"categoryname"`
	Username     string `json:"username"`
}

func (cgr *ValueControlRequest) Bind(r *http.Request) error {
	if cgr == nil {
		return errors.New("missing required user fields")
	}

	return nil
}
