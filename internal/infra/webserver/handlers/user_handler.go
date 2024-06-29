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

type Error struct {
	Message string `json:"message"`
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

func (handler *UserHandler) GetJwt(w http.ResponseWriter, req *http.Request) {
	jwt := req.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := req.Context().Value("jwtExpiresIn").(int64)
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

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
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

// Create user godoc
// @Summary 		Create user
// @Description 	Create user
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			request	body	dto.CreateUserInput	true	"user request"
// @Success 		201
// @Failure 		500 	{object}	Error
// @Failure 		400 	{object}	Error
// @Router 			/users 	[post]
func (handler *UserHandler) CreateUser(w http.ResponseWriter, req *http.Request) {
	var userInput dto.CreateUserInput
	err := json.NewDecoder(req.Body).Decode(&userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		apiError := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(apiError)
		return
	}

	u, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		apiError := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(apiError)
		return
	}

	err = handler.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		apiError := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(apiError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
