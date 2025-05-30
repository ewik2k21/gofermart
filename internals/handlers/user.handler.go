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

// FillDb
// @Summary Заполняет базу данных случайными пользователями
// @Description Заполняет таблицу users в базе данных указанным количеством случайных пользователей.
// @Tags users
// @Accept json
// @Produce json
// @Param count body int true "Количество пользователей для добавления"
// @Success 200 {string} string "Successfully fill db"
// @Failure 400 {object} interfaces.Response "Invalid request"
// @Failure 500 {object} interfaces.Response "Internal Server Error"
// @Router /user/filldb [post]
func (h *UserHandler) FillDb(c *gin.Context) {
	var count int
	if err := c.ShouldBindJSON(&count); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := h.userService.FillDb(count, c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, "Successfully fill db")

}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user account and generates a JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param user body interfaces.UserRequest true "User registration details"
// @Success 200 {object} interfaces.Response{Code=int, Message=string} "Successful registration"
// @Failure 400 {object} interfaces.Response{Code=int, Message=string} "Invalid request or validation error"
// @Failure 500 {object} interfaces.Response{Code=int, Message=string} "Internal server error"
// @Router /user/register [post]
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

// Login godoc
// @Summary Logs in a user
// @Description Authenticates a user and returns a JWT token.
// @Tags users
// @Accept json
// @Produce json
// @Param user body interfaces.UserRequest true "User credentials"
// @Success 200 {object} interfaces.Response{Message=string} "User successfully logged in"
// @Failure 400 {object} interfaces.Response{Message=string} "Bad Request"
// @Failure 401 {object} interfaces.Response{Message=string} "Unauthorized"
// @Failure 500 {object} interfaces.Response{Message=string} "Internal Server Error"
// @Router /user/login [post]
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

// AddOrder godoc
// @Summary Adds a new order for the user
// @Description Adds a new order to the user's account based on the provided order number.
// @Tags users, orders
// @Accept json
// @Produce json
// @Param order body interfaces.OrderRequest true "Order details"
// @Security ApiKeyAuth
// @Success 200 {object} interfaces.Response{Message=string} "Successfully added order"
// @Success 202 {object} interfaces.Response{Message=string} "Order accepted for processing"
// @Failure 400 {object} interfaces.Response{Message=string} "Bad Request"
// @Failure 401 {object} interfaces.Response{Message=string} "Unauthorized"
// @Failure 422 {object} interfaces.Response{Message=string} "Unprocessable Entity"
// @Failure 500 {object} interfaces.Response{Message=string} "Internal Server Error"
// @Router /user/orders [post]
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
		return
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

// GetAllOrders godoc
// @Summary Gets all orders for the user
// @Description Retrieves a list of all orders associated with the authenticated user.
// @Tags users, orders
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} interfaces.OrderResponse "Successfully retrieved orders"
// @Failure 204 {object} interfaces.Response{Message=string} "No Content"
// @Failure 401 {object} interfaces.Response{Message=string} "Unauthorized"
// @Failure 500 {object} interfaces.Response{Message=string} "Internal Server Error"
// @Router /user/orders [get]
func (h *UserHandler) GetAllOrders(c *gin.Context) {
	userId, ok := utils.GetId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Wrong user id",
			Code:    http.StatusUnauthorized,
		})
		return
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

// GetBalance godoc
// @Summary Gets the user's balance
// @Description Retrieves the current balance for the authenticated user.
// @Tags users, balance
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {number} float64 "Successfully retrieved balance"
// @Failure 401 {object} interfaces.Response{Message=string} "Unauthorized"
// @Failure 500 {object} interfaces.Response{Message=string} "Internal Server Error"
// @Router /user/balance [get]
func (h *UserHandler) GetBalance(c *gin.Context) {
	userId, ok := utils.GetId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Wrong user id",
			Code:    http.StatusUnauthorized,
		})
	}

	balance, err := h.userService.GetBalance(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: "Failed get balance",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// GetWithdraws godoc
// @Summary Get user withdraws
// @Description Retrieves all withdraws for a specific user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} interfaces.WithdrawsResponse "Successfully retrieved withdraws"
// @Failure 401 {object} interfaces.Response{code=int,message=string} "Unauthorized (Wrong user ID)"
// @Failure 204 {object} interfaces.Response{code=int,message=string} "No Content (No withdraws found)"
// @Failure 500 {object} interfaces.Response{code=int,message=string} "Internal Server Error"
// @Router /user/withdraws [get]
func (h *UserHandler) GetWithdraws(c *gin.Context) {
	userId, ok := utils.GetId(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, interfaces.Response{
			Message: "Wrong user id",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	withdraws, err := h.userService.GetWithdrawsById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	if len(*withdraws) == 0 {
		c.JSON(http.StatusNoContent, interfaces.Response{
			Message: "No withdraws",
			Code:    http.StatusNoContent,
		})
	}

	c.JSON(http.StatusOK, withdraws)
}

// PostWithdraw godoc
// @Summary Post a withdraw request
// @Description Endpoint to submit a withdraw request
// @Tags User
// @Accept json
// @Produce json
// @Param request body interfaces.WithdrawRequest true "Withdraw Request"
// @Security ApiKeyAuth
// @Success 200 {object} interfaces.Response{code=int,message=string} "Success"
// @Failure 400 {object} interfaces.Response{code=int,message=string} "Bad Request (Invalid JSON, validation errors)"
// @Failure 422 {object} interfaces.Response{code=int,message=string} "Unprocessable Entity (Incorrect order number format)"
// @Failure 401 {object} interfaces.Response{code=int,message=string} "Unauthorized (Wrong user ID)"
// @Failure 500 {object} interfaces.Response{code=int,message=string} "Internal Server Error"
// @Router /user/balance/withdraw [post]
func (h *UserHandler) PostWithdraw(c *gin.Context) {
	var withdrawRequest interfaces.WithdrawRequest

	if err := c.ShouldBindJSON(&withdrawRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.validate.Struct(withdrawRequest); err != nil {
		c.JSON(http.StatusBadRequest, interfaces.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if ok := utils.ValidateOrderNumber(withdrawRequest.Order); !ok {
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
		return
	}

	statusCode, msg, err := h.userService.PostWithdraw(userId, withdrawRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, interfaces.Response{
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(statusCode, interfaces.Response{Message: msg, Code: statusCode})
}
