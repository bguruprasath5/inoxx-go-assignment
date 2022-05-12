package models

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	UserID   uint   `json:"id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}
