package token

import (
	"os"

	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func Init() {
	secretPhrase := os.Getenv("SECRET_PASSPHRASE")

	TokenAuth = jwtauth.New("HS256", []byte(secretPhrase), nil)
}
