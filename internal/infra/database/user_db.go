package database

import (
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{DB: db}
}

func (udb *UserDB) Create(user *entity.User) error {
	return udb.DB.Create(user).Error
}

func (udb *UserDB) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := udb.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
