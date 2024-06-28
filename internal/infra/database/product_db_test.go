package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := CreateGormDBAndAutoMigrate()
	assert.NoError(t, err)

	product, err := entity.NewProduct("Pimponeta", 10000.0)
	assert.NoError(t, err)

	prodDB := NewProduct(db)
	err = prodDB.Create(product)
	assert.NoError(t, err)

	var productFound *entity.Product
	err = db.Find(&productFound, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestFindProductByID(t *testing.T) {
	db, err := CreateGormDBAndAutoMigrate()
	assert.NoError(t, err)

	product, err := entity.NewProduct("TestFindProductByID", 1000.0)
	assert.NoError(t, err)

	err = db.Create(product).Error
	assert.NoError(t, err)

	prodDB := NewProduct(db)
	productFound, err := prodDB.FindByID(product.ID.String())
	assert.NoError(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := CreateGormDBAndAutoMigrate()
	assert.NoError(t, err)

	product, err := entity.NewProduct("TestUpdateProduct", 1000.0)
	assert.NoError(t, err)

	err = db.Create(product).Error
	assert.NoError(t, err)

	var createdProductFound *entity.Product
	err = db.First(&createdProductFound, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, product.ID, createdProductFound.ID)
	assert.Equal(t, product.Name, createdProductFound.Name)
	assert.Equal(t, product.Price, createdProductFound.Price)

	product.Name = "UpdatedProduct"
	product.Price = 100

	prodDB := NewProduct(db)
	prodDB.Update(product)

	var updatedProductFound *entity.Product
	err = db.First(&updatedProductFound, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, updatedProductFound.Name, "UpdatedProduct")
	assert.Equal(t, updatedProductFound.Price, 100.0)
}

func TestDeleteProduct(t *testing.T) {
	db, err := CreateGormDBAndAutoMigrate()
	assert.NoError(t, err)

	product, err := entity.NewProduct("TestDeleteProduct", 1.0)
	assert.NoError(t, err)

	err = db.Create(product).Error
	assert.NoError(t, err)

	prodDB := NewProduct(db)
	var foundProduct *entity.Product
	err = db.First(&foundProduct, "id = ?", product.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, product.ID, foundProduct.ID)
	assert.Equal(t, product.Name, foundProduct.Name)
	assert.Equal(t, product.Price, foundProduct.Price)

	err = prodDB.Delete(foundProduct.ID.String())
	assert.NoError(t, err)

	err = db.First(&foundProduct, "id = ?", product.ID).Error
	assert.Error(t, err)
}

func TestFindAllProducts(t *testing.T) {
	db, err := CreateGormDBAndAutoMigrate()
	assert.NoError(t, err)

	for i := 1; i < 24; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("Product %d", i), rand.Float64()*100)
		assert.NoError(t, err)
		err = db.Create(product).Error
		assert.NoError(t, err)
	}

	productDB := NewProduct(db)
	products, err := productDB.FindAll(1, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = productDB.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = productDB.FindAll(1, 10, "desc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 23", products[0].Name)
	assert.Equal(t, "Product 14", products[9].Name)

	products, err = productDB.FindAll(2, 10, "desc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 13", products[0].Name)
	assert.Equal(t, "Product 4", products[9].Name)
}

func CreateGormDBAndAutoMigrate() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entity.Product{})
	return db, nil
}
