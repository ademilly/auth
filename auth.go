package auth

import (
	"encoding/json"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"golang.org/x/crypto/bcrypt"
)

// User struct holds user data
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hash     string `json:"hash,omitempty"`
}

// UserClosure should enclose User retrieval information
type UserClosure func() (User, error)

// Hash hashes plaintext password of user
func Hash(user User) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, err
	}
	return User{
		Username: user.Username,
		Hash:     string(hash),
	}, nil
}

// CheckHash compares plain text candidate password to hashed password obtained from getRegistered closure
func CheckHash(candidate User, getRegistered UserClosure) error {
	registered, err := getRegistered()
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(registered.Hash), []byte(candidate.Password))
}

// Protect ensures `handler` is behind tokenMiddleware and requires valid token to be used
func Protect(tokenMiddleware *jwtmiddleware.JWTMiddleware, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenMiddleware.Handler(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

// UserFromRequest reads request body into a User struct
func UserFromRequest(r *http.Request) (User, error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var candidate User
	err := decoder.Decode(&candidate)

	if err != nil {
		return User{}, err
	}
	return candidate, nil
}

// Login builds a http.HandleFunc handling login request
func Login(tokenizer func(string) (string, error), getRegistered UserClosure) func(User) (string, error) {
	return func(candidate User) (string, error) {
		err := CheckHash(candidate, getRegistered)
		if err != nil {
			return "", err
		}

		return tokenizer(candidate.Username)
	}
}
