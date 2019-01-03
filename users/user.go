package users

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-bongo/bongo"
	"gitlab.com/benCoder01/payword-backend/db"
	"gitlab.com/benCoder01/payword-backend/mail"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/sethvargo/go-password/password"
	"gitlab.com/benCoder01/payword-backend/requests"
	"gitlab.com/benCoder01/payword-backend/responses"
	"gitlab.com/benCoder01/payword-backend/token"
)

// Router verwaltet alle Routes, die mit /users beginnen.
func Router(cors *cors.Cors) chi.Router {
	r := chi.NewRouter()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler)

		r.Post("/sign-in", signIn)
		r.Post("/sign-up", signUp)
		r.Post("/mail/reset-password", setNewPassword)
	})

	// Protected
	r.Group(func(r chi.Router) {
		r.Use(cors.Handler)

		r.Use(jwtauth.Verifier(token.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/change-password", changePassword)
		r.Post("/mail/update-mail", saveEmail)
	})

	/*
		r.Get("/", getUsers)
		r.Get("/{name}", getUser)

		r.Get("/increment/{name}/{category}", increment)
		r.Get("/decrement/{name}/{category}", decrement)
	*/

	return r
}

func signIn(w http.ResponseWriter, r *http.Request) {
	userReq := &requests.UserRequest{}

	if err := render.Bind(r, userReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	user, err := db.FindUserByName(userReq.Username)

	if err != nil {

		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrInvalidRequest(errors.New("Could not validate user")))
			return
		}

		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if !compareHashToPassword(user.Password, userReq.Password) {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("Could not validate user")))
		return
	}

	//Token
	_, token, err := token.TokenAuth.Encode(jwt.MapClaims{"username": user.Username})

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		returndocker run --name payword-backend -p 27017:27017 -d mongo

	}

	if err := render.Render(w, r, responses.NewTokenResponse(token)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}

}

func signUp(w http.ResponseWriter, r *http.Request) {
	userReq := &requests.UserRequest{}

	if err := render.Bind(r, userReq); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}

	exists, err := db.UserExists(userReq.Username)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if exists {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New("User already exists")))
		return
	}

	password, err := createHash(userReq.Password)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	user := &db.User{Username: userReq.Username, Password: password}

	err = db.SaveUser(user)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	pwchangeReq := &requests.PasswordChangeRequest{}

	if err := render.Bind(r, pwchangeReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	user, err := db.FindUserByName(pwchangeReq.Username)

	if err != nil {

		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrInternal(errors.New("Could not validate user")))
			return
		}

		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if !compareHashToPassword(user.Password, pwchangeReq.OldPassword) {
		render.Render(w, r, responses.ErrInternal(errors.New("Could not validate user")))
		return
	}

	newHashedPassword, err := createHash(pwchangeReq.NewPassword)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	user.Password = newHashedPassword

	err = db.SaveUser(user)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func createHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func compareHashToPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		return false
	}

	return true
}

func saveEmail(w http.ResponseWriter, r *http.Request) {
	emailReq := &requests.EmailControlRequest{}

	if err := render.Bind(r, emailReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	user, err := db.FindUserByName(emailReq.Username)

	if err != nil {
		if db.NotFoundError(err) {
			render.Render(w, r, responses.ErrInvalidRequest(errors.New("User not found")))
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	// If a user already saved his mail, the current adress will be replaced
	mail, err := db.GetMailAdress(emailReq.Username)

	if err != nil {
		if db.NotFoundError(err) {
			mail = &db.Mail{Username: emailReq.Username, Mail: emailReq.Mail}
		} else {
			render.Render(w, r, responses.ErrInternal(err))
			return
		}
	} else {
		mail.Mail = emailReq.Mail
	}

	if err := render.Render(w, r, responses.NewMailResponse(mail.Username, mail.Mail)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func setNewPassword(w http.ResponseWriter, r *http.Request) {
	emailReq := &requests.EmailControlRequest{}

	if err := render.Bind(r, emailReq); err != nil {
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	mailInfo, err := db.GetMailAdress(emailReq.Username)

	if err != nil {
		if db.NotFoundError(err) {
			render.Render(w, r, responses.ErrInvalidRequest(errors.New("Invalid credentials")))
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	user, err := db.FindUserByName(mailInfo.Username)

	if err != nil {
		if db.NotFoundError(err) {
			render.Render(w, r, responses.ErrInvalidRequest(errors.New("Invalid credentials")))
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	generatedPassword, err := generateNewPassword()

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	hashedPassword, err := createHash(generatedPassword)

	if err != nil {
		render.Render(w, r, responses.ErrInternal(err))
		return
	}

	user.Password = hashedPassword

	err = db.SaveUser(user)

	err = mail.Send(mailInfo.Mail, hashedPassword)

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func overwritePassword(user *db.User, password string) error {
	user.Password = password
	return db.SaveUser(user)
}

func generateNewPassword() (string, error) {
	// TODO: generate shorter password
	return password.Generate(64, 10, 10, false, false)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := db.FindAllUsers()
	if err := render.RenderList(w, r, responses.NewUserListResponse(users)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "name")

	user, err := db.FindUserByName(username)

	if err != nil {

		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			render.Render(w, r, responses.ErrNotFound())
		} else {
			render.Render(w, r, responses.ErrInternal(err))
		}
		return
	}

	if err := render.Render(w, r, responses.NewUserResponse(user)); err != nil {
		render.Render(w, r, responses.ErrRender(err))
		return
	}
}
