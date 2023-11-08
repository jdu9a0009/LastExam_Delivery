package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"api-gateway-service/config"
	"api-gateway-service/genproto/order_service"
	user_service "api-gateway-service/genproto/user_service"
	"api-gateway-service/pkg/helper"

	"github.com/gin-gonic/gin"
)

// Task3 branchlarni activeni hozirgi vaqtga nisbatlab olish

// GetAllBranch godoc
// @Security ApiKeyAuth
// @Router       /v1/branch/active [get]
// @Summary      GetAll Active branch
// @Description  get allactive branch
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Success      200  {array}   user_service.ListBranchResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListActiveBranch(ctx *gin.Context) {
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
	timeNow := time.Now().Format("15:04:05")

	resp, err := h.services.BranchService().ListActive(ctx.Request.Context(), &user_service.ListActiveBranchRequest{
		Page:  int32(page),
		Limit: int32(limit),
		Date:  timeNow,
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListActiveBranch", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllBranch response", http.StatusOK, resp)
}

//task 5 va 6)

// 5. Zakaz statusini o'zgartirish uchun endpoint(API)

// Updated Order Status godoc
// @Security ApiKeyAuth
// @Router       /v1/logic [get]
// @Summary      Update order status
// @Description  Update order status
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        order_id    query     int  false  "order_id for response"
// @Success      200  {string}   string
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// Get the previous status
	prevStatus, err := h.services.OrderService().GetOrderStatus(ctx.Request.Context(), &order_service.OrderIdRequest{
		OrderId: idStr,
	})
	if err != nil {
		h.handlerResponse(ctx, "error get prevstatus", http.StatusBadRequest, err.Error())
		return
	}

	// Extract the status value from the previous status response
	status := prevStatus.Status

	// Update the status
	response, err := h.services.OrderService().UpdateStatus(ctx.Request.Context(), &order_service.UpdateOrderStatusRequest{
		OrderId: idStr,
		Status:  status,
	})

	// Handle the update response and error (if any)
	if err != nil {
		h.handlerResponse(ctx, "failed to update order status", http.StatusBadRequest, err.Error())
		return
	}

	// Return the response to the client
	h.handlerResponse(ctx, "order status updated successfully", http.StatusOK, response)

	// 6. Zakaz zavershit bo'lganda clientni ma'lumotlari(last_ordered_date,total_orders_count,total_orders_sum)ni update qilish

	if status == "finished" {
		// Get Order using orderid
		respOrder, err := h.services.OrderService().Get(ctx.Request.Context(), &order_service.IdStrRequest{
			Id: idStr,
		})
		if err != nil {
			h.handlerResponse(ctx, "error get order in updateOrderStatus", http.StatusBadRequest, err.Error())
			return
		}

		_, err = h.services.ClientService().UpdateOrder(ctx.Request.Context(), &user_service.UpdateClientsOrderRequest{
			Id:               respOrder.ClientId,
			TotalOrdersCount: 1,
			TotalOrdersSum:   respOrder.Price,
		})
		if err != nil {
			h.handlerResponse(ctx, "error update client order in updateOrderStatus", http.StatusBadRequest, err.Error())
			return
		}
	}
}

//Task 7. Courier endpoints[get]:
//   7.1 Qabul qilishi mumkin bo'lgan zakazlar(courier olmagan,statuslari mos kelgan)
//   7.2 Qabul qilgan va yakunlanmagan zakazlar

// GetListOrder of Courier godoc
// @Router       /v1/courier/get_order/{id} [get]
// @Summary      Get a courier by ID
// @Description  Get accepted and not accepted orders using courier_id
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Courier ID to retrieve"
// @Success      200  {object}  order_service.Order
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetCourierOrders(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error courier GetById", http.StatusBadRequest, err.Error())
		return
	}
	respAcceptedOrders, err := h.services.OrderService().GetAllAcceptedOrders(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error courier GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "getAll Accepted Orders of this courier response", http.StatusOK, respAcceptedOrders)
	respAcceptableOrders, err := h.services.OrderService().GetAllAcceptableOrders(ctx.Request.Context(), &order_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error courier GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "getAll AcceptableOrders  of this courier response", http.StatusOK, respAcceptableOrders)
}

// 9.Zakazdan courierni olib tashlash uchun endpoint:
//  - zakazni statusi 'Accepted'ga o'zgaradi

// @Router       /v1/logic/{id} [get]
// @Summary      Zakazda courierni olib tashlash
// @Description  api for update order
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of order"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) DeleteCourierInOrder(c *gin.Context) {
	idStr := c.Param("id")

	ID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(c, "error order update order status parse id", http.StatusBadRequest, err.Error())
		return
	}

	respOrder, err := h.services.OrderService().Get(c.Request.Context(), &order_service.IdStrRequest{Id: idStr})
	if err != nil {
		h.handlerResponse(c, "error courier GetById", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.OrderService().Update(c.Request.Context(), &order_service.UpdateOrderRequest{

		Id:            int32(ID),
		OrderId:       respOrder.OrderId,
		ClientId:      respOrder.ClientId,
		BranchId:      respOrder.BranchId,
		CourierId:     0,
		Type:          respOrder.Type,
		Address:       respOrder.Address,
		DeliveryPrice: respOrder.DeliveryPrice,
		Price:         respOrder.Discount,
		Discount:      respOrder.Discount,
		PaymentType:   respOrder.PaymentType,
		Status:        "accepted",
	})
	if err != nil {
		h.handlerResponse(c, "error courier Update", http.StatusInternalServerError, err.Error())
		return
	}
	h.handlerResponse(c, "DeletedCourirInOrder", http.StatusOK, resp)

}

// 	8. Courier zakazni olishi uchun endpiont(API):
//  - zakazda courier bor bo'lsa error qaytarish,
//  - courierda max orders countga teng zakazlari bo'lsa error qaytarish
//  - zakaz statusi 'Courier Accepted'ga o'zgaradi

// @Router       /v1/courier/delete_order/{id} [get]
// @Summary      Update Order
// @Description  courier get Order
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        id    path     string  true  "id of order"
// @Success      200  {string}   string
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp

func (h *Handler) CourierGetOrder(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")

	token, err := helper.ExtractToken(tokenStr)
	if err != nil {
		h.handlerResponse(c, "error get courier", http.StatusBadRequest, err.Error())
		return
	}

	claims, err := helper.ParseClaims(token, config.JWTSecretKey)
	if err != nil {
		h.handlerResponse(c, "error get claims", http.StatusBadRequest, err.Error())
		return
	}

	respCourier, err := h.services.CourierService().Get(c.Request.Context(), &user_service.IdRequest{Id: claims.UserID})
	if err != nil {
		h.handlerResponse(c, "error get Couriers", http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")

	respOrder, err := h.services.OrderService().Get(c.Request.Context(), &order_service.IdStrRequest{Id: id})

	if respOrder.CourierId != 0 {
		c.JSON(http.StatusBadRequest, "order already received")
		return
	}
	idstr := int32(respOrder.CourierId)
	respCouriersOrder, err := h.services.OrderService().GetAllAcceptableOrders(c.Request.Context(), &order_service.IdRequest{
		Id: idstr})

	if err != nil {
		fmt.Println("error Order Get:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	if respCouriersOrder.Count == respCourier.MaxOrderCount {
		c.JSON(http.StatusBadRequest, "you are reached max order!")
		return
	}

	resp, err := h.services.OrderService().UpdateStatus(c.Request.Context(), &order_service.UpdateOrderStatusRequest{
		OrderId: respOrder.OrderId,
		Status:  "courier_accepted",
	})

	if err != nil {
		fmt.Println("error Order Update:", err.Error())
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	h.handlerResponse(c, "Courier Get Order", http.StatusOK, resp)

}
