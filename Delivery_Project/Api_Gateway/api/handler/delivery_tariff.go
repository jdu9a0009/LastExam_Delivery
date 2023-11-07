package handler

import (
	"fmt"
	"net/http"
	"strconv"

	order_service "api-gateway-service/genproto/order_service"

	"github.com/gin-gonic/gin"
)

// CreateDeliveryTariff godoc
// @Router       /v1/delivery_tariff [post]
// @Summary      Create a new delivery_tariff
// @Description  Create a new delivery_tariff with the provided details
// @Tags         delivery_tariff
// @Accept       json
// @Produce      json
// @Param        delivery_tariff     body  order_service.CreateDeliveryTariffRequest true  "data of the delivery_tariff"
// @Success      201  {object}  order_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateDeliveryTariff(ctx *gin.Context) {
	var delivery_tariff = order_service.CreateDeliveryTariffRequest{}

	err := ctx.ShouldBindJSON(&delivery_tariff)
	if err != nil {
		h.handlerResponse(ctx, "CreateDeliveryTariff", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.DeliveryTariffService().Create(ctx, &order_service.CreateDeliveryTariffRequest{
		Name:       delivery_tariff.Name,
		TariffType: delivery_tariff.TariffType,
		BasePrice:  delivery_tariff.BasePrice,
		Values:     delivery_tariff.Values,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create delivery_tariff response", http.StatusOK, resp)
}

// GetAllDeliveryTariff godoc
// @Security ApiKeyAuth
// @Router       /v1/delivery_tariff [get]
// @Summary      GetAll DeliveryTariff
// @Description  get delivery_tariff
// @Tags         delivery_tariff
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        name     query     string false "search by name"
// @Param        tariff_type     query     string false "search by tariff_type"
// @Success      200  {array}   order_service.ListDeliveryTariffResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListDeliveryTariff(ctx *gin.Context) {
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

	resp, err := h.services.DeliveryTariffService().List(ctx.Request.Context(), &order_service.ListDeliveryTariffRequest{
		Page:      int32(page),
		Limit:     int32(limit),
		Search:    ctx.Query("search"),
		TarifType: ctx.Query("tariff_type"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListDeliveryTariff", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllDeliveryTariff response", http.StatusOK, resp)
}

// GetDeliveryTariff godoc
// @Security ApiKeyAuth
// @Router       /v1/delivery_tariff/{id} [get]
// @Summary      Get a delivery_tariff by ID
// @Description  Retrieve a delivery_tariff by its unique identifier
// @Tags         delivery_tariff
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "DeliveryTariff ID to retrieve"
// @Success      200  {object}  order_service.DeliveryTariff
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetDeliveryTariff(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.DeliveryTariffService().Get(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error delivery_tariff GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get delivery_tariff response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/delivery_tariff/{id} [put]
// @Summary      Update an existing delivery_tariff
// @Description  Update an existing delivery_tariff with the provided details
// @Tags         delivery_tariff
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "DeliveryTariff ID to update"
// @Param        delivery_tariff   body    order_service.UpdateDeliveryTariffRequest  true    "Updated data for the delivery_tariff"
// @Success      200  {object}  order_service.UpdateDeliveryTariffRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateDeliveryTariff(ctx *gin.Context) {
	var delivery_tariff = order_service.UpdateDeliveryTariffRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&delivery_tariff)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	delivery_tariff.Id = int32(id)

	resp, err := h.services.DeliveryTariffService().Update(ctx.Request.Context(), &delivery_tariff)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error order Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update order response", http.StatusOK, resp)
}

// DeleteDeliveryTariff godoc
// @Security ApiKeyAuth
// @Router       /v1/delivery_tariff/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a delivery_tariff by its unique identifier
// @Tags         delivery_tariff
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  order_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteDeliveryTariff(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.DeliveryTariffService().Delete(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error delivery_tariff Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete delivery_tariff response", http.StatusOK, resp)
}
