package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/NeGat1FF/todolist-api/internal/service"
	"github.com/NeGat1FF/todolist-api/internal/utils"
)

type RegisterUserRequest struct {
	Username string
	Email    string
	Password string
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type TokensResponse struct {
	Token        string
	RefreshToken string
}

type RefreshTokenResponse struct {
	RefreshToken string
}

type UserHandler struct {
	ser *service.Service
}

func NewUserHandler(ser *service.Service) *UserHandler {
	return &UserHandler{ser}
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body			RegisterUserRequest	true	"User object that needs to be registered"
//	@Success		200		{object}	TokensResponse
//	@Router			/users/register [post]
func (uh *UserHandler) RegisterUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(models.UserKey{}).(models.User)

	tokenString, refreshTokenString, err := uh.ser.RegisterUser(r.Context(), user)
	if err != nil {
		http.Error(rw, err.Error(), err.(service.ServerError).Code)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]any{"token": tokenString, "refreshToken": refreshTokenString})
}

// LoginUser godoc
//
//	@Summary		Login a user
//	@Description	Login a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body			LoginUserRequest	true	"User object that needs to be logged in"
//	@Success		200		{object}	TokensResponse
//	@Router			/users/login [post]
func (uh *UserHandler) LoginUser(rw http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(models.UserKey{}).(models.User)

	tokenString, refreshTokenString, err := uh.ser.LoginUser(r.Context(), user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusExpectationFailed)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]any{"token": tokenString, "refreshToken": refreshTokenString})
}

// RefreshToken godoc
//
//	@Summary		Refresh a token
//	@Description	Refresh a token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200		{object}	RefreshTokenResponse
//	@Router			/users/refresh [post]
func (uh *UserHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	tokenString := strings.TrimSpace(strings.Replace(r.Header.Get("Authorization"), "Bearer", "", -1))
	if tokenString == "" {
		http.Error(rw, "no authorization token", http.StatusUnauthorized)
		return
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		http.Error(rw, "failed to validate token", http.StatusUnauthorized)
		return
	}

	var user_id int

	expTime := claims["exp"].(float64)
	if time.Now().After(time.Unix(int64(expTime), 0)) {
		http.Error(rw, "token expired", http.StatusUnauthorized)
		return
	}
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		http.Error(rw, "invalid type of token", http.StatusUnauthorized)
		return
	}
	user_id = int(claims["uid"].(float64))

	newToken := uh.ser.IssueAccessToken(user_id)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]any{"token": newToken})
}
