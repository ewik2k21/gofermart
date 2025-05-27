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

//func (h *UserHandler) GetId(c *gin.Context) {
//	userId, exists := c.Get("user_id")
//	if !exists {
//		c.JSON(400, "FAILED ID ")
//		return
//	}
//	userIdString, _ := userId.(string)
//	c.JSON(200, "user id : "+userIdString)
//}
