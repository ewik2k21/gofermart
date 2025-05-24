package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type RouteDefinition struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

type Claims struct {
	UserId string
	jwt.StandardClaims
}
