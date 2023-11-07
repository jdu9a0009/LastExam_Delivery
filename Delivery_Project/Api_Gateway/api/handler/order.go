package handler

import (
	"fmt"
	"net/http"
	"strconv"

	order_service "api-gateway-service/genproto/order_service"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
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

	err := ctx.ShouldBindJSON(&order)
	if err != nil {
		h.handlerResponse(ctx, "CreateOrder", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.OrderService().Create(ctx, &order_service.CreateOrderRequest{
		ClientId:      order.ClientId,
		BranchId:      order.BranchId,
		Type:          order.Type,
		CourierId:     order.CourierId,
		DeliveryPrice: order.DeliveryPrice,
		Price:         order.Price,
		Discount:      order.Discount,
		PaymentType:   order.PaymentType,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create order response", http.StatusOK, resp)
}

// GetAllOrder godoc
// @Router       /v1/order [get]
// @Summary      GetAll Order
// @Description  get order
// @Tags         order
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        name     query     string false "search by title"
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
	/*
			    BranchId     int32   `protobuf:"varint,5,opt,name=branch_id,json=branchId,proto3" json:"branch_id,omitempty"`
		    DeliveryType string  `protobuf:"bytes,6,opt,name=delivery_type,json=deliveryType,proto3" json:"delivery_type,omitempty"`
		    CourierId    int32   `protobuf:"varint,7,opt,name=courier_id,json=courierId,proto3" json:"courier_id,omitempty"`
		    PriceFrom    float64 `protobuf:"fixed64,8,opt,name=price_from,json=priceFrom,proto3" json:"price_from,omitempty"`
		    PriceTo      float64 `protobuf:"fixed64,9,opt,name=price_to,json=priceTo,proto3" json:"price_to,omitempty"`
		    PaymentType  string  `protobuf:"bytes,10,opt,name=payment_type,json=paymentType,proto3" json:"payment_type,omitempty"`
		}

	*/
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
	id := ctx.Param("id")

	resp, err := h.services.OrderService().Get(ctx.Request.Context(), &order_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error order GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get order response", http.StatusOK, resp)
}

// UpdateOrder godoc
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
	id := ctx.Param("id")

	resp, err := h.services.OrderService().Delete(ctx.Request.Context(), &order_service.IdRequest{Id: id})
	if err != nil {
		h.handlerResponse(ctx, "error order Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete order response", http.StatusOK, resp)
}
