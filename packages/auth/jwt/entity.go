package jwtAuth

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtData struct {
	UserID    int64  `json:"user_id"`
	Reference string `json:"reference"`
}

type JwtPayload struct {
	Reference string `json:"reference"`
	UserID    int64  `json:"ui"`
	jwt.StandardClaims
}
