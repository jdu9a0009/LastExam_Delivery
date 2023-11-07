package handler

import (
	"fmt"
	"net/http"
	"strconv"

	user_service "api-gateway-service/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// CreateClients godoc
// @Router       /v1/client [post]
// @Summary      Create a new client
// @Description  Create a new client with the provided details
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        client     body  user_service.CreateClientsRequest true  "data of the client"
// @Success      201  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateClients(ctx *gin.Context) {
	var client = user_service.CreateClientsRequest{}

	err := ctx.ShouldBindJSON(&client)
	if err != nil {
		h.handlerResponse(ctx, "CreateClients", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.ClientService().Create(ctx, &user_service.CreateClientsRequest{
		Firstname:      client.Firstname,
		Lastname:       client.Lastname,
		Phone:          client.Phone,
		Photo:          client.Photo,
		BirthDate:      client.BirthDate,
		DiscountType:   client.DiscountType,
		DiscountAmount: client.DiscountAmount,
	})

	if err != nil {
		h.handlerResponse(ctx, "ClientsService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create client response", http.StatusOK, resp)
}

// GetAllClients godoc
// @Security ApiKeyAuth
// @Router       /v1/client [get]
// @Summary      GetAll Clients
// @Description  get client
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        search     query     string false "search by lastnam,firstname and phone"
// @Param        created_at_from     query     string false "search by created_at_from"
// @Param        created_at_to     query     string false "search by created_at_to"
// @Param        last_order_date_from     query     string false "search by last_order_date_from"
// @Param        last_order_date_to    query     string false "search by last_order_date_to"
// @Param        total_orders_count_from    query     int false "search by total_orders_count_from"
// @Param        total_orders_count_to    query     int false "search by total_orders_sum_to"
// @Param        total_orders_sum_from    query     int false "search by total_orders_sum_from"
// @Param        total_orders_sum_to    query     int false "search by total_orders_count_to"
// @Param        discount_type    query     string false "search by discount_type"
// @Param        discount_from    query     string false "search by discount_from"
// @Param        discount_to    query     string false "search by discount_to"
// @Success      200  {array}   user_service.ListClientsResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListClients(ctx *gin.Context) {
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

	total_orders_count_from, err := h.ParseQueryParam(ctx, "total_orders_count_from", "0")
	if err != nil {
		fmt.Println(err)
		return
	}

	total_orders_count_to, err := h.ParseQueryParam(ctx, "total_orders_count_to", "0")
	if err != nil {
		fmt.Println(err)
		return
	}

	total_orders_sum_from, err := h.ParseQueryParam(ctx, "total_orders_sum_from", "0")
	if err != nil {
		fmt.Println(err)
		return
	}
	total_orders_sum_to, err := h.ParseQueryParam(ctx, "total_orders_sum_to", "0")
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := h.services.ClientService().List(ctx.Request.Context(), &user_service.ListClientsRequest{
		Page:                 int32(page),
		Limit:                int32(limit),
		Search:               ctx.Query("search"),
		CreatedAtFrom:        ctx.Query("created_at_from"),
		CreatedAtTo:          ctx.Query("created_at_to"),
		LastOrderedDateFrom:  ctx.Query("last_order_date_from"),
		LastOrderedDateTo:    ctx.Query("last_order_date_to"),
		TotalOrdersCountFrom: int64(total_orders_count_from),
		TotalOrdersCountTo:   int64(total_orders_count_to),
		TotalOrdersSumFrom:   int64(total_orders_sum_from),
		TotalOrdersSumTo:     int64(total_orders_sum_to),
		DiscountType:         ctx.Query("dicount_type"),
		DiscountAmountFrom:   ctx.Query("discount_amount_from"),
		DiscountAmountTo:     ctx.Query("discount_amount_to"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListClients", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllClients response", http.StatusOK, resp)
}

// GetClients godoc
// @Security ApiKeyAuth
// @Router       /v1/client/{id} [get]
// @Summary      Get a client by ID
// @Description  Retrieve a client by its unique identifier
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Clients ID to retrieve"
// @Success      200  {object}  user_service.Clients
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetClients(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.ClientService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error client GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get client response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/client/{id} [put]
// @Summary      Update an existing client
// @Description  Update an existing client with the provided details
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Clients ID to update"
// @Param        client   body    user_service.UpdateClientsRequest  true    "Updated data for the client"
// @Success      200  {object}  user_service.UpdateClientsRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateClients(ctx *gin.Context) {
	var client = user_service.UpdateClientsRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&client)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	client.Id = int32(id)

	resp, err := h.services.ClientService().Update(ctx.Request.Context(), &client)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error user Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update user response", http.StatusOK, resp)
}

// DeleteClients godoc
// @Security ApiKeyAuth
// @Router       /v1/client/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a client by its unique identifier
// @Tags         client
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteClients(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.ClientService().Delete(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error client Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete client response", http.StatusOK, resp)
}
