package models

import "github.com/golang-jwt/jwt/v5"

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys represents a public key.
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Claims struct {
	Email string `json:"email"`
	Token_type string `json:"type"`
	jwt.RegisteredClaims
}