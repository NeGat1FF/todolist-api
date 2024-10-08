package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/NeGat1FF/todolist-api/internal/models"
)

func validateUser(r *http.Request, checkUsername bool) (models.User, error) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		return user, errors.New("failed to parse request body")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	switch {
	case checkUsername && user.Username == "":
		return user, errors.New("username is not specified")
	case user.Password == "":
		return user, errors.New("password is not specified")
	case user.Email == "":
		return user, errors.New("email is not specified")
	case !regex.MatchString(user.Email):
		return user, errors.New("invalid email address")
	}

	return user, nil
}

func ValidateRegistration(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, err := validateUser(r, true)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserKey{}, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

func ValidateLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		user, err := validateUser(r, false)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserKey{}, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}
