package handler

import (
	"fmt"
	"net/http"
	"strconv"

	product_service "api-gateway-service/genproto/product_service"

	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/product [post]
// @Summary      Create a new product
// @Description  Create a new product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product     body  product_service.CreateProductRequest true  "data of the product"
// @Success      201  {object}  product_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateProduct(ctx *gin.Context) {
	var product = product_service.CreateProductRequest{}

	err := ctx.ShouldBindJSON(&product)
	if err != nil {
		h.handlerResponse(ctx, "CreateProduct", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.ProductService().Create(ctx, &product_service.CreateProductRequest{
		Title:       product.Title,
		Description: product.Description,
		Photo:       product.Photo,
		OrderNumber: product.OrderNumber,
		ProductType: product.ProductType,
		Price:       product.Price,
		CategoryId:  product.CategoryId,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create product response", http.StatusOK, resp)
}

// GetAllProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/product [get]
// @Summary      GetAll Product
// @Description  get product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param search query string false "search by title"
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param type query integer false "type"
// @Param category_id query int false "category_id"
// @Success      200  {array}   product_service.ListProductResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListProduct(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		h.handlerResponse(ctx, "error get page", http.StatusBadRequest, err.Error())
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		h.handlerResponse(ctx, "error get limit", http.StatusBadRequest, err.Error())
		return
	}
	category_id, err := h.ParseQueryParam(ctx, "category_id", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := h.services.ProductService().List(ctx.Request.Context(), &product_service.ListProductRequest{
		Page:     int32(page),
		Limit:    int32(limit),
		Search:   ctx.Query("title"),
		Type:     ctx.Query("type"),
		Category: int32(category_id),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListProduct", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllProduct response", http.StatusOK, resp)
}

// GetProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/product/{id} [get]
// @Summary      Get a product by ID
// @Description  Retrieve a product by its unique identifier
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Product ID to retrieve"
// @Success      200  {object}  product_service.Product
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.services.ProductService().Get(ctx.Request.Context(), &product_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error product GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get product response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/product/{id} [put]
// @Summary      Update an existing product
// @Description  Update an existing product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Product ID to update"
// @Param        product   body    product_service.UpdateProductRequest  true    "Updated data for the product"
// @Success      200  {object}  product_service.UpdateProductRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateProduct(ctx *gin.Context) {
	var product = product_service.UpdateProductRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&product)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	product.Id = int32(id)

	resp, err := h.services.ProductService().Update(ctx.Request.Context(), &product)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error product Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update product response", http.StatusOK, resp)
}

// DeleteProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/product/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a product by its unique identifier
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  product_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.services.ProductService().Delete(ctx.Request.Context(), &product_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error product Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete product response", http.StatusOK, resp)
}
