package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	p, err := NewProduct("Product 1", 100.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.Equal(t, "Product 1", p.Name)
	assert.Equal(t, 100.0, p.Price)
	assert.NotNil(t, p.CreatedAt)
}

func TestNewProductWhenNameIsRequired(t *testing.T) {
	p, err := NewProduct("", 100.0)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestNewProductWhenPriceIsRequired(t *testing.T) {
	p, err := NewProduct("Product 1", 0.0)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestNewProductWhenPriceIsInvalid(t *testing.T) {
	p, err := NewProduct("Product 1", -1.0)
	assert.NotNil(t, err)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestProductValidate(t *testing.T) {
	p, err := NewProduct("Product 1", 100.0)
	assert.Nil(t, err)
	assert.Nil(t, p.Validate())
}
