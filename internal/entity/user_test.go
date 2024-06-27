package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("Mr. Pipo", "chandelier.pipo@gmail.com", "123321")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "Mr. Pipo", user.Name)
	assert.Equal(t, "chandelier.pipo@gmail.com", user.Email)
	assert.NotEqual(t, "123321", user.Password)
	assert.NotEmpty(t, user.Password)
}

func TestValidatePassword(t *testing.T) {
	user, err := NewUser("Mr. Pipo", "chandelier.pipo@gmail.com", "123321")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123321"))
	assert.False(t, user.ValidatePassword("1233212"))
	assert.NotEqual(t, "123321", user.Password)
}
