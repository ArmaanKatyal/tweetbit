package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email      string `json:"email"`
	Token_type string `json:"type"`
	jwt.RegisteredClaims
}
