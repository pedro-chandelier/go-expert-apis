package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/pedro-chandelier/go-expert-apis/configs"
	_ "github.com/pedro-chandelier/go-expert-apis/docs"
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/database"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title 			Go Expert API
// @version 		1.0
// @description 	Product API with authentication

// @contact.name 	Pedro Haeser
// @contact.email 	chandelier.pipo@gmail.com

// @host 							localhost:8000
// @BasePath 						/
// @securityDefinitions.apiKey 		ApiKeyAuth
// @in 								header
// @name 							Authorization
func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	configs := configs.LoadConfig("configs/.env")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.WithValue("jwt", configs.TokenAuth))
	router.Use(middleware.WithValue("jwtExpiresIn", configs.JwtExpiresIn))

	router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	attachUserHandler(db, router)
	attachProductHandler(db, router)

	http.ListenAndServe(":8000", router)
}

func attachProductHandler(db *gorm.DB, router *chi.Mux) {
	configs := configs.LoadConfig("configs/.env")
	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)

	router.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

}

func attachUserHandler(db *gorm.DB, router *chi.Mux) {
	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB)

	router.Post("/users", userHandler.CreateUser)
	router.Post("/users/generate-token", userHandler.GetJwt)
}
