package handler

import (
	"fmt"
	"net/http"
	"strconv"

	order_service "api-gateway-service/genproto/order_service"
	"api-gateway-service/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Security ApiKeyAuth
// @Router       /v1/order [post]
// @Summary      Create a new order
// @Description  Create a new order with the provided details
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        order     body  order_service.CreateOrderRequest true  "data of the order"
// @Success      201  {object}  order_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateOrder(ctx *gin.Context) {
	var order = order_service.CreateOrderRequest{}
	var discountprice float32
	var deliveryprice float32
	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		h.handlerResponse(ctx, "CreateOrder", http.StatusBadRequest, err.Error())
		return
	}

	respClient, err := h.services.ClientService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: order.ClientId})
	if err != nil {
		h.handlerResponse(ctx, "error client GetById in create order", http.StatusBadRequest, err.Error())
		return
	}

	if respClient.DiscountType == "percent" {
		discountprice = float32(order.Price) * float32(respClient.DiscountAmount)
	}
	if respClient.DiscountType == "fixed" {
		discountprice = float32(respClient.DiscountAmount)
	}

	//task 10 delvery priceni hisoblash
	respBranch, err := h.services.BranchService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: order.BranchId})
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}

	respDeliveryTariff, err := h.services.DeliveryTariffService().Get(ctx.Request.Context(), &order_service.IdRequest{Id: respBranch.DeliveryTarifId})
	if err != nil {
		h.handlerResponse(ctx, "error delivery_tariff GetById in order create", http.StatusBadRequest, err.Error())
		return
	}

	if respDeliveryTariff.TariffType == "fixed" {
		deliveryprice = float32(respDeliveryTariff.BasePrice)
	} else if respDeliveryTariff.TariffType == "alternative" {
		respTariff, err := h.services.DeliveryTariffService().List(ctx.Request.Context(), &order_service.ListDeliveryTariffRequest{
			Page:      1,
			Limit:     10,
			TarifType: "alternative",
		})

		if err != nil {
			h.handlerResponse(ctx, "error GetListDeliveryTariff in  create order", http.StatusBadRequest, err.Error())
			return
		}
		for _, v := range respTariff.DeliveryTariffs {
			if order.Price > v.Values.FromPrice && order.Price < v.Values.ToPrice {
				deliveryprice = float32(v.Values.Price)

			}
		}
	}

	resp, err := h.services.OrderService().Create(ctx, &order_service.CreateOrderRequest{
		ClientId:      order.ClientId,
		BranchId:      order.BranchId,
		Type:          order.Type,
		CourierId:     order.CourierId,
		DeliveryPrice: float64(deliveryprice),
		Price:         order.Price - float64(discountprice),
		Discount:      float64(discountprice),
		PaymentType:   order.PaymentType,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create order response", http.StatusOK, resp)
}

// GetAllOrder godoc
// @Security ApiKeyAuth
// @Router       /v1/order [get]
// @Summary      GetAll Order
// @Description  get order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param search query string false "search"
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param page query integer false "page"
// @Param order_id query integer false "order_id"
// @Param client_id query integer false "client_id"
// @Param branch_id query integer false "branch_id"
// @Param delivery_type query string false "delivery_type"
// @Param courier_id query integer false "courier_id"
// @Param price_from query integer false "price_from"
// @Param price_to query integer false "price_to"
// @Param payment_type query string false "payment_type"
// @Success      200  {array}   order_service.ListOrderResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListOrder(ctx *gin.Context) {
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
	client_id, err := h.ParseQueryParam(ctx, "client_id", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	branch_id, err := h.ParseQueryParam(ctx, "branch_id", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	courier_id, err := h.ParseQueryParam(ctx, "courier_id", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	price_from, err := h.ParseQueryParam(ctx, "price_from", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	price_to, err := h.ParseQueryParam(ctx, "price_to", "0")
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := h.services.OrderService().List(ctx.Request.Context(), &order_service.ListOrderRequest{
		Limit:        int32(limit),
		Page:         int32(page),
		OrderId:      ctx.Query("order_id"),
		ClientId:     int32(client_id),
		BranchId:     int32(branch_id),
		CourierId:    int32(courier_id),
		DeliveryType: ctx.Query("delivery_type"),
		PriceFrom:    float64(price_from),
		PriceTo:      float64(price_to),
		PaymentType:  ctx.Query("payment_type"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListOrder", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllOrder response", http.StatusOK, resp)
}

// GetOrder godoc
// @Security ApiKeyAuth
// @Router       /v1/order/{id} [get]
// @Summary      Get a order by ID
// @Description  Retrieve a order by its unique identifier
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Order ID to retrieve"
// @Success      200  {object}  order_service.Order
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.OrderService().Get(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error order GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get order response", http.StatusOK, resp)
}

// UpdateOrder godoc
// @Security ApiKeyAuth
// @Router       /v1/order/{id} [put]
// @Summary      Update an existing order
// @Description  Update an existing order with the provided details
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Order ID to update"
// @Param        order   body    order_service.UpdateOrderRequest  true    "Updated data for the order"
// @Success      200  {object}  order_service.UpdateOrderRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateOrder(ctx *gin.Context) {
	var order = order_service.UpdateOrderRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&order)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	order.Id = int32(id)

	resp, err := h.services.OrderService().Update(ctx.Request.Context(), &order)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error order Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update order response", http.StatusOK, resp)
}

// DeleteOrder godoc
// @Security ApiKeyAuth
// @Router       /v1/order/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a order by its unique identifier
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  order_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.OrderService().Delete(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error order Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete order response", http.StatusOK, resp)
}
