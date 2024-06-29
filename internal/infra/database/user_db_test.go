package database

import (
	"testing"

	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory::"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, err := entity.NewUser("Mr. Pipo", "chandelier.pipo@gmail.com", "pipolino")
	if err != nil {
		t.Error(err)
	}

	userDB := NewUserDB(db)
	err = userDB.Create(user)
	assert.Nil(t, err)

	var userFound *entity.User
	err = db.Find(&userFound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, user.ID, userFound.ID)
	assert.Equal(t, "chandelier.pipo@gmail.com", userFound.Email)
	assert.Equal(t, "Mr. Pipo", userFound.Name)
	assert.NotEqual(t, "pipolino", userFound.Password)
	assert.NotEmpty(t, userFound.Password)
}

func TestUserFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory::"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, err := entity.NewUser("Mr. Pipo", "chandelier.pipo@gmail.com", "pipolino")
	if err != nil {
		t.Error(err)
	}

	userDB := NewUserDB(db)
	err = userDB.Create(user)
	assert.Nil(t, err)

	userFound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	assert.NotNil(t, userFound)
	assert.Equal(t, user.Email, userFound.Email)
	assert.Equal(t, user.Name, userFound.Name)
	assert.NotNil(t, userFound.Password)
	assert.NotEmpty(t, userFound.Password)
}
