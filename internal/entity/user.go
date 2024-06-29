package entity

import (
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/validator"
	"github.com/pedro-chandelier/go-expert-apis/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required"`
	Password string    `json:"-" validate:"required"`
}

func NewUser(name, email, password string) (*User, error) {
	err := validator.GetValidatorInstance().Struct(&User{
		Name:     name,
		Email:    email,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
