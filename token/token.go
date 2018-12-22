package token

import (
	"os"

	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func Init() {
	secretPhrase := os.Getenv("SECRET_PASSPHRASE")
	// fmt.Println(secretPhrase)
	TokenAuth = jwtauth.New("HS256", []byte(secretPhrase), nil)
}
