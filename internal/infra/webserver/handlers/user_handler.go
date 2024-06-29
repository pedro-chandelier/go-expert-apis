package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/pedro-chandelier/go-expert-apis/internal/dto"
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/database"
)

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(userDB database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDB:       userDB,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

func (handler *UserHandler) GetJwt(w http.ResponseWriter, req *http.Request) {
	var jwtInput dto.GetJwtInput
	err := json.NewDecoder(req.Body).Decode(&jwtInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := handler.UserDB.FindByEmail(jwtInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !user.ValidatePassword(jwtInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := handler.Jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(handler.JwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	var userInput dto.CreateUserInput
	err := json.NewDecoder(req.Body).Decode(&userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = handler.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
