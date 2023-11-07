package handler

import (
	"fmt"
	"net/http"
	"strconv"

	user_service "api-gateway-service/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// CreateBranch godoc
// @Router       /v1/branch [post]
// @Summary      Create a new branch
// @Description  Create a new branch with the provided details
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        branch     body  user_service.CreateBranchRequest true  "data of the branch"
// @Success      201  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) CreateBranch(ctx *gin.Context) {
	var branch = user_service.CreateBranchRequest{}

	err := ctx.ShouldBindJSON(&branch)
	if err != nil {
		h.handlerResponse(ctx, "CreateBranch", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.BranchService().Create(ctx, &user_service.CreateBranchRequest{
		Name:            branch.Name,
		Photo:           branch.Photo,
		Phone:           branch.Phone,
		DeliveryTarifId: branch.DeliveryTarifId,
		WorkHourStart:   branch.WorkHourStart,
		WorkHourEnd:     branch.WorkHourEnd,
		Address:         branch.Address,
		Destination:     branch.Destination,
	})

	if err != nil {
		h.handlerResponse(ctx, "BranchService().Create", http.StatusBadRequest, err.Error())

		return
	}

	h.handlerResponse(ctx, "create branch response", http.StatusOK, resp)
}

// GetAllBranch godoc
// @Security ApiKeyAuth
// @Router       /v1/branch [get]
// @Summary      GetAll Branch
// @Description  get branch
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        limit    query     int  false  "limit for response"  Default(10)
// @Param		 page     query     int  false  "page for response"   Default(1)
// @Param        name     query     string false "search by name"
// @Param        created_at_from     query     string false "search by created_at_from"
// @Param        created_at_to     query     string false "search by created_at_to"
// @Success      200  {array}   user_service.ListBranchResponse
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetListBranch(ctx *gin.Context) {
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

	resp, err := h.services.BranchService().List(ctx.Request.Context(), &user_service.ListBranchRequest{
		Page:          int32(page),
		Limit:         int32(limit),
		Name:          ctx.Query("name"),
		CreatedAtFrom: ctx.Query("created_at_from"),
		CreatedAtTo:   ctx.Query("created_at_to"),
	})

	if err != nil {
		h.handlerResponse(ctx, "error GetListBranch", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get AllBranch response", http.StatusOK, resp)
}

// GetBranch godoc
// @Security ApiKeyAuth
// @Router       /v1/branch/{id} [get]
// @Summary      Get a branch by ID
// @Description  Retrieve a branch by its unique identifier
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id   path    string     true    "Branch ID to retrieve"
// @Success      200  {object}  user_service.Branch
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) GetBranch(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}
	resp, err := h.services.BranchService().Get(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error branch GetById", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "get branch response", http.StatusOK, resp)
}

// UpdateProduct godoc
// @Security ApiKeyAuth
// @Router       /v1/branch/{id} [put]
// @Summary      Update an existing branch
// @Description  Update an existing branch with the provided details
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id       path    int     true    "Branch ID to update"
// @Param        branch   body    user_service.UpdateBranchRequest  true    "Updated data for the branch"
// @Success      200  {object}  user_service.UpdateBranchRequest
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) UpdateBranch(ctx *gin.Context) {
	var branch = user_service.UpdateBranchRequest{}
	id, err := strconv.Atoi(ctx.Param("id"))

	err = ctx.ShouldBind(&branch)
	if err != nil {
		h.handlerResponse(ctx, "error while binding", http.StatusBadRequest, err.Error())
		return
	}

	branch.Id = int32(id)

	resp, err := h.services.BranchService().Update(ctx.Request.Context(), &branch)
	fmt.Println("before  send bind", resp)
	if err != nil {
		h.handlerResponse(ctx, "error user Update", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "update user response", http.StatusOK, resp)
}

// DeleteBranch godoc
// @Security ApiKeyAuth
// @Router       /v1/branch/{id} [delete]
// @Summary      Delete a Catgory
// @Description  delete a branch by its unique identifier
// @Tags         branch
// @Accept       json
// @Produce      json
// @Param        id   path    int     true    "Catgeory ID to retrieve"
// @Success      200  {object}  user_service.Response
// @Failure      400  {object}  Response{data=string}
// @Failure      404  {object}  Response{data=string}
// @Failure      500  {object}  Response{data=string}
func (h *Handler) DeleteBranch(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.handlerResponse(ctx, "error branch Delete", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.services.BranchService().Delete(ctx.Request.Context(), &user_service.IdRequest{Id: int32(id)})
	if err != nil {
		h.handlerResponse(ctx, "error branch Delete", http.StatusBadRequest, err.Error())
		return
	}

	h.handlerResponse(ctx, "delete branch response", http.StatusOK, resp)
}
