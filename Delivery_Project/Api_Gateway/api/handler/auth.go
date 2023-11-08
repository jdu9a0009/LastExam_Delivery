package handler

import (
	"api-gateway-service/api/response"
	"api-gateway-service/config"
	user_service "api-gateway-service/genproto/user_service"
	"api-gateway-service/pkg/helper"
	"api-gateway-service/pkg/logger"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRes struct {
	Token string `json:"token"`
}

// @Router       /login [post]
// @Summary      create user
// @Description  api for create users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user    body     LoginReq  true  "data of user"
// @Success      200  {object}  LoginRes
// @Failure      400  {object}  response.ErrorResp
// @Failure      404  {object}  response.ErrorResp
// @Failure      500  {object}  response.ErrorResp
func (h *Handler) Login(c *gin.Context) {

	// Task 2

	var req LoginReq

	err := c.ShouldBindJSON(&req)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "BAD REQUEST", Message: "invalid fields in body"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	hashPass, err := helper.GeneratePasswordHash(req.Password)

	if err != nil {
		h.log.Error("error while binding:", logger.Error(err))
		res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if req.Role == "courier" {
		resp, err := h.services.CourierService().GetCourierByUserName(c.Request.Context(), &user_service.GetByUserName{
			Login: req.Login,
		})

		if err != nil {
			fmt.Println("error Staff GetByLoging:", err.Error())
			res := response.ErrorResp{Code: "INTERNAL ERROR", Message: "internal server error"}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		err = helper.ComparePasswords([]byte(hashPass), []byte(resp.Password))

		if err != nil {
			h.log.Error("error while binding:", logger.Error(err))
			res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		m := make(map[string]interface{})
		m["user_id"] = resp.Id
		token, err := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)

		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, LoginRes{Token: token})

		return

	} else if req.Role == "user" {

		resp, err := h.services.UserService().GetUserByUserName(c.Request.Context(), &user_service.GetByUserName{
			Login: req.Login,
		})

		if err != nil {
			fmt.Println("error Staff GetByLoging:", err.Error())
			res := response.ErrorResp{Code: "INTERNAL ERROR", Message: "internal server error"}
			c.JSON(http.StatusInternalServerError, res)
			return
		}

		err = helper.ComparePasswords([]byte(hashPass), []byte(resp.Password))

		if err != nil {
			h.log.Error("error while binding:", logger.Error(err))
			res := response.ErrorResp{Code: "INVALID Password", Message: "invalid password"}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		m := make(map[string]interface{})
		m["user_id"] = resp.Id
		token, err := helper.GenerateJWT(m, config.TokenExpireTime, config.JWTSecretKey)

		if err != nil {
			return
		}

		c.JSON(http.StatusCreated, LoginRes{Token: token})

		return

	} else {
		fmt.Println("Role undefined")
		res := response.ErrorResp{Code: "INVALID Role", Message: "invalid role"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

}
