package main

import (
	"net/http"

	"github.com/pedro-chandelier/go-expert-apis/configs"
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/database"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig("configs/.env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}
