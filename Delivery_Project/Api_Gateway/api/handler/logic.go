package handler

import (
	"net/http"
	"strconv"
	"time"

	"api-gateway-service/genproto/order_service"
	user_service "api-gateway-service/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// Task3 branchlarni activeni hozirgi vaqtga nisbatlab olish

// GetAllActiveBranch godoc
// @Security ApiKeyAuth
// @Router       /v1/logic [get]
// @Summary      GetAll Active Branch
// @Description  get branch
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

//t(ask 5 va 6)

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
// @Success      200  {array}   user_service.ListOrderResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateOrderStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")

	orderID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error order update order status parse id", http.StatusBadRequest, err.Error())
		return
	}

	// Get the previous status
	prevStatus, err := h.services.OrderService().GetOrderStatus(ctx.Request.Context(), &order_service.OrderIdRequest{
		OrderId: int32(orderID),
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
		respOrder, err := h.services.OrderService().Get(ctx.Request.Context(), &order_service.IdRequest{
			Id: int32(orderID),
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
// @Router       /v1/logic/{id} [get]
// @Summary      Get a courier by ID
// @Description  Get accepted and not accepted orders using courier_id
// @Tags         logic
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Courier ID to retrieve"
// @Success      200  {object}  order_service.Orders
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
