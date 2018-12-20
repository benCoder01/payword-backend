package token

import "github.com/go-chi/jwtauth"

var TokenAuth *jwtauth.JWTAuth

func Init() {
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}
