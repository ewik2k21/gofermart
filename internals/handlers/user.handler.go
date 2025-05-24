package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gofermart/internals/interfaces"
	"gofermart/internals/services"
	"net/http"
	"time"
)

type UserHandler struct {
	userService  services.UserService
	tokenService services.TokenService
	validate     *validator.Validate
}

func NewUserHandler(userService services.UserService, tokenService services.TokenService) *UserHandler {
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

	c.Header("Authorization", "Bearer"+*tokenString)
	c.Header("Token-Expiration", expirationTime.Format(time.RFC3339))

	c.JSON(http.StatusOK, interfaces.Response{
		Message: "User create successfully and logged in",
		Code:    http.StatusOK,
	})
}
