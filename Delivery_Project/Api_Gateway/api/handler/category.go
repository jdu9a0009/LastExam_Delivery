package handler

import (
	"fmt"
	"net/http"
	"strconv"

	product_service "api-gateway-service/genproto/product_service"

	"github.com/gin-gonic/gin"
)

// CreateCategory godoc
// @Router       /v1/category [post]
// @Summary      Create a new category
// @Description  Create a new category with the provided details
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        category     body  product_service.CreateCategoryRequest true  "data of the category"
// @Success      201  {object}  product_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateCategory(ctx *gin.Context) {
	var category = product_service.CreateCategoryRequest{}

	err := ctx.ShouldBindJSON(&category)
	if err != nil {
		h.handlerResponse(ctx, "CreateCategory", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.CategoryService().Create(ctx, &product_service.CreateCategoryRequest{
		Title:       category.Title,
		Image:       category.Image,
		ParentId:    category.ParentId,
		OrderNumber: category.OrderNumber,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create category response", http.StatusOK, resp)
}

// GetAllCategory godoc
// @Router       /v1/category [get]
// @Summary      GetAll Category
// @Description  get category
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        title     query     string false "search by title"
// @Success      200  {array}   product_service.ListCategoryResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListCategory(ctx *gin.Context) {
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

	resp, err := h.services.CategoryService().List(ctx.Request.Context(), &product_service.ListCategoryRequest{
		Page:   int32(page),
		Limit:  int32(limit),
		Search: ctx.Query("title"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListCategory", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllCategory response", http.StatusOK, resp)
}

// GetCategory godoc
// @Router       /v1/category/{id} [get]
// @Summary      Get a category by ID
// @Description  Retrieve a category by its unique identifier
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Category ID to retrieve"
// @Success      200  {object}  product_service.Category
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.services.CategoryService().Get(ctx.Request.Context(), &product_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error category GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get category response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Router       /v1/category/{id} [put]
// @Summary      Update an existing category
// @Description  Update an existing category with the provided details
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Category ID to update"
// @Param        category   body    product_service.UpdateCategoryRequest  true    "Updated data for the category"
// @Success      200  {object}  product_service.UpdateCategoryRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateCategory(ctx *gin.Context) {
	var category = product_service.UpdateCategoryRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&category)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	category.Id = int32(id)

	resp, err := h.services.CategoryService().Update(ctx.Request.Context(), &category)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error product Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update product response", http.StatusOK, resp)
}

// DeleteCategory godoc
// @Router       /v1/category/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a category by its unique identifier
// @Tags         category
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  product_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := h.services.CategoryService().Delete(ctx.Request.Context(), &product_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error category Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete category response", http.StatusOK, resp)
}
