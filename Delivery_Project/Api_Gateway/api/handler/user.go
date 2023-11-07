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

// CreateUser godoc
// @Router       /v1/user [post]
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user     body  user_service.CreateUsersRequest true  "data of the user"
// @Success      201  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateUser(ctx *gin.Context) {
	var user = user_service.CreateUsersRequest{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		h.handlerResponse(ctx, "CreateUser", http.StatusBadRequest, err.Error())
		return
	}
	hashedPass, err := helper.GeneratePasswordHash(user.Password)
	if err != nil {
		h.log.Error("error while generating hash password:", logger.Error(err))
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}
	user.Password = string(hashedPass)

	resp, err := h.services.UserService().Create(ctx, &user_service.CreateUsersRequest{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		BranchId:  user.BranchId,
		Phone:     user.Phone,
		Login:     user.Login,
		Password:  user.Password,
	})

	if err != nil {
		h.handlerResponse(ctx, "CatgeoryService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create user response", http.StatusOK, resp)
}

// GetAllUser godoc
// @Security ApiKeyAuth
// @Router       /v1/user [get]
// @Summary      GetAll User
// @Description  get user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        name     query     string false "search by firstname,lastname and phone"
// @Param        created_at_from     query     string false "search by created_at_from"
// @Param        created_at_to     query     string false "search by created_at_to"
// @Success      200  {array}   user_service.ListUsersResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListUser(ctx *gin.Context) {
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

	resp, err := h.services.UserService().List(ctx.Request.Context(), &user_service.ListUsersRequest{
		Page:          int32(page),
		Limit:         int32(limit),
		Search:        ctx.Query("title"),
		CreatedAtFrom: ctx.Query("created_at_from"),
		CreatedAtTo:   ctx.Query("created_at_to"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListUser", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllUser response", http.StatusOK, resp)
}

// GetUser godoc
// @Security ApiKeyAuth
// @Router       /v1/user/{id} [get]
// @Summary      Get a user by ID
// @Description  Retrieve a user by its unique identifier
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "User ID to retrieve"
// @Success      200  {object}  user_service.Users
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.UserService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error user GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get user response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/user/{id} [put]
// @Summary      Update an existing user
// @Description  Update an existing user with the provided details
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "User ID to update"
// @Param        user   body    user_service.UpdateUsersRequest  true    "Updated data for the user"
// @Success      200  {object}  user_service.UpdateUsersRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateUser(ctx *gin.Context) {
	var user = user_service.UpdateUsersRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&user)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	user.Id = int32(id)

	resp, err := h.services.UserService().Update(ctx.Request.Context(), &user)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error user Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update user response", http.StatusOK, resp)
}

// DeleteUser godoc
// @Security ApiKeyAuth
// @Router       /v1/user/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a user by its unique identifier
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Catgeory ID to retrieve"
// @Success      200  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.UserService().Delete(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error user Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete user response", http.StatusOK, resp)
}
