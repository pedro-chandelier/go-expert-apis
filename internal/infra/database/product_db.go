package database

import (
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"gorm.io/gorm"
)

type ProductDB struct {
	DB *gorm.DB
}

func NewProduct(db *gorm.DB) *ProductDB {
	return &ProductDB{DB: db}
}

func (pdb *ProductDB) Create(product *entity.Product) error {
	return pdb.DB.Create(product).Error
}

func (pdb *ProductDB) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := pdb.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (pdb *ProductDB) Update(product *entity.Product) error {
	_, err := pdb.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return pdb.DB.Save(product).Error
}

func (pdb *ProductDB) Delete(id string) error {
	product, err := pdb.FindByID(id)
	if err != nil {
		return err
	}
	return pdb.DB.Delete(product).Error
}

func (pdb *ProductDB) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = pdb.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		err = pdb.DB.Order("created_at " + sort).Find(&products).Error
	}
	return products, err
}
