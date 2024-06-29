package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/products", productHandler.CreateProduct)
	router.Get("/products", productHandler.GetProducts)
	router.Get("/products/{id}", productHandler.GetProduct)
	router.Put("/products/{id}", productHandler.UpdateProduct)
	router.Delete("/products/{id}", productHandler.DeleteProduct)

	http.ListenAndServe(":8000", router)
}
