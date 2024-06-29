package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pedro-chandelier/go-expert-apis/internal/dto"
	"github.com/pedro-chandelier/go-expert-apis/internal/entity"
	"github.com/pedro-chandelier/go-expert-apis/internal/infra/database"
	entityPkg "github.com/pedro-chandelier/go-expert-apis/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
// @Summary 		Create product
// @Description 	Create product
// @Tags 			products
// @Accept 			json
// @Produce 		json
// @Param 			request		body	dto.CreateProductInput	true	"product request"
// @Success 		201
// @Failure 		400 		{object}	Error
// @Failure 		500 		{object}	Error
// @Router 			/products 	[post]
// @Security		ApiKeyAuth
func (handler *ProductHandler) CreateProduct(w http.ResponseWriter, req *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{Message: err.Error()})
		return
	}

	err = handler.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetProducts godoc
// @Summary 		Get all products
// @Description 	Get all products
// @Tags 			products
// @Accept 			json
// @Produce 		json
// @Param 			page		query	string	false	"page number"
// @Param 			page		query	string	false	"limit"
// @Success 		200			{array}	entity.Product
// @Failure 		404 		{object}	Error
// @Failure 		500 		{object}	Error
// @Router 			/products 	[get]
// @Security		ApiKeyAuth
func (handler *ProductHandler) GetProducts(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Query().Get("page")
	limit := req.URL.Query().Get("limit")
	sort := req.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := handler.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// GetProduct godoc
// @Summary 		Get a product
// @Description 	Get a product
// @Tags 			products
// @Accept 			json
// @Produce 		json
// @Param 			id				path		string		true 	"product ID"	Format(uuid)
// @Success 		200				{object}	entity.Product
// @Failure 		404
// @Failure 		500 			{object}	Error
// @Router 			/products/{id} 	[get]
// @Security		ApiKeyAuth
func (handler *ProductHandler) GetProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := handler.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary 		Update a product
// @Description 	Update a product
// @Tags 			products
// @Accept 			json
// @Produce 		json
// @Param 			id				path		string		true 	"product ID"	Format(uuid)
// @Param 			request			body		dto.CreateProductInput		true 	"product request"
// @Success 		200
// @Failure 		404
// @Failure 		500
// @Router 			/products/{id} 	[put]
// @Security		ApiKeyAuth
func (handler *ProductHandler) UpdateProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = handler.ProductDB.FindByID(product.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = handler.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary 		Delete a product
// @Description 	Delete a product
// @Tags 			products
// @Accept 			json
// @Produce 		json
// @Param 			id				path		string		true 	"product ID"	Format(uuid)
// @Success 		200
// @Failure 		404
// @Failure 		500
// @Router 			/products/{id} 	[delete]
// @Security		ApiKeyAuth
func (handler *ProductHandler) DeleteProduct(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := handler.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = handler.ProductDB.Delete(product.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
