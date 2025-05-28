package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gofermart/internals/interfaces"
	"gofermart/internals/services"
	"gofermart/internals/utils"
	"net/http"
	"time"
)

type UserHandler struct {
	userService  services.IUserService
	tokenService services.ITokenService
	validate     *validator.Validate
}

func NewUserHandler(userService services.IUserService, tokenService services.ITokenService) *UserHandler {
	return &UserHandler{userService: userService, tokenService: tokenService, validate: validator.New()}
}

func (h *UserHandler) Register(c *gin.Context) {
	var userRequest interfaces.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	if err := h.validate.Struct(userRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	userId, err := h.userService.CreateUserAccount(&userRequest, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	tokenString, expirationTime, err := h.tokenService.GenerateJwtToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.Header("Authorization", "Bearer "+*tokenString)
	c.Header("Token-Expiration", expirationTime.Format(time.RFC3339))

	c.JSON(http.StatusOK, interfaces.Response{
		Message: "User create successfully and logged in",
		Code:    http.StatusOK,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var userRequest interfaces.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.validate.Struct(userRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})

		return
	}

	userId, passwordOk, err := h.userService.CheckCredentials(&userRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Login or password incorrect",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	if !passwordOk {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Login or password incorrect",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	tokenString, expirationTime, err := h.tokenService.GenerateJwtToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.Header("Authorization", "Bearer "+*tokenString)
	c.Header("Token-Expiration", expirationTime.Format(time.RFC3339))

	c.JSON(http.StatusOK, interfaces.Response{
		Message: "User successfully logged in",
		Code:    http.StatusOK,
	})

}

func (h *UserHandler) AddOrder(c *gin.Context) {
	var orderRequest interfaces.OrderRequest

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.validate.Struct(orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if ok := utils.ValidateOrderNumber(orderRequest.OrderNumber); !ok {
		c.JSON(http.StatusUnprocessableEntity, interfaces.Response{
			Code:    http.StatusUnprocessableEntity,
			Message: "Incorrect order number format ",
		})
		return
	}
	userId, ok := utils.GetId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Wrong user id",
			Code:    http.StatusUnauthorized,
		})
	}
	statusCode, message, err := h.userService.AddOrder(userId, orderRequest.OrderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(statusCode, interfaces.Response{
		Message: message,
		Code:    statusCode,
	})
}

func (h *UserHandler) GetAllOrders(c *gin.Context) {
	userId, ok := utils.GetId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Wrong user id",
			Code:    http.StatusUnauthorized,
		})
	}

	orders, err := h.userService.GetAllOrders(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	if len(*orders) == 0 {
		c.JSON(http.StatusNoContent, interfaces.Response{
			Message: "No data to answer",
			Code:    http.StatusNoContent,
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

//func (h *UserHandler) GetId(c *gin.Context) {
//	userId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(400, "FAILED ID ")
//		return
//	}
//	userIdString, _ := userId.(string)
//	c.JSON(200, "user id : "+userIdString)
//}
