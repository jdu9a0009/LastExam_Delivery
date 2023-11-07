package handler

import (
	"fmt"
	"net/http"
	"strconv"

	user_service "api-gateway-service/genproto/user_service"
	"api-gateway-service/pkg/helper"
	"api-gateway-service/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CreateCourier godoc
// @Router       /v1/courier [post]
// @Summary      Create a new courier
// @Description  Create a new courier with the provided details
// @Tags         courier
// @Accept       json
// @Produce      json
// @Param        courier     body  user_service.CreateCouriersRequest true  "data of the courier"
// @Success      201  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateCourier(ctx *gin.Context) {
	var courier = user_service.CreateCouriersRequest{}

	err := ctx.ShouldBindJSON(&courier)
	if err != nil {
		h.handlerResponse(ctx, "CreateCourier", http.StatusBadRequest, err.Error())
		return
	}

	hashedPass, err := helper.GeneratePasswordHash(courier.Password)
	if err != nil {
		h.log.Error("error while generating hash password:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	courier.Password = string(hashedPass)

	resp, err := h.services.CourierService().Create(ctx, &user_service.CreateCouriersRequest{
		Firstname:     courier.Firstname,
		Lastname:      courier.Lastname,
		BranchId:      courier.BranchId,
		Phone:         courier.Phone,
		Login:         courier.Login,
		Password:      courier.Password,
		MaxOrderCount: courier.MaxOrderCount,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create courier response", http.StatusOK, resp)
}

// GetAllCourier godoc
// @Security ApiKeyAuth
// @Router       /v1/courier [get]
// @Summary      GetAll Courier
// @Description  get courier
// @Tags         courier
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        name     query     string false "Search by firstname and lastname and phone"
// @Param        created_at_from     query     string false "search by created_at_from"
// @Param        created_at_to     query     string false "search by created_at_to"
// @Success      200  {array}   user_service.ListCouriersResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListCourier(ctx *gin.Context) {
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

	resp, err := h.services.CourierService().List(ctx.Request.Context(), &user_service.ListCouriersRequest{
		Page:          int32(page),
		Limit:         int32(limit),
		Search:        ctx.Query("title"),
		CreatedAtFrom: ctx.Query("created_at_from"),
		CreatedAtTo:   ctx.Query("created_at_to"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListCourier", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllCourier response", http.StatusOK, resp)
}

// GetCourier godoc
// @Security ApiKeyAuth
// @Router       /v1/courier/{id} [get]
// @Summary      Get a courier by ID
// @Description  Retrieve a courier by its unique identifier
// @Tags         courier
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Courier ID to retrieve"
// @Success      200  {object}  user_service.Couriers
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetCourier(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.CourierService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error courier GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get courier response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/courier/{id} [put]
// @Summary      Update an existing courier
// @Description  Update an existing courier with the provided details
// @Tags         courier
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Courier ID to update"
// @Param        courier   body    user_service.UpdateCouriersRequest  true    "Updated data for the courier"
// @Success      200  {object}  user_service.UpdateCouriersRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateCourier(ctx *gin.Context) {
	var courier = user_service.UpdateCouriersRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&courier)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	courier.Id = int32(id)

	resp, err := h.services.CourierService().Update(ctx.Request.Context(), &courier)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error courier Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update courier response", http.StatusOK, resp)
}

// DeleteCourier godoc
// @Security ApiKeyAuth
// @Router       /v1/courier/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a courier by its unique identifier
// @Tags         courier
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteCourier(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.CourierService().Delete(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error courier Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete courier response", http.StatusOK, resp)
}
