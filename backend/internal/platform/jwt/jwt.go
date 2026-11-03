package jwt

import (
	"fmt"
	"time"

	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secretKey     string
	expirationHrs int
}

type Claims struct {
	PostulanteID int    `json:"postulante_id"`
	Cedula       string `json:"cedula"`
	jwt.RegisteredClaims
}

func NewTokenManager(secretKey string, expirationHrs int) *TokenManager {
	return &TokenManager{
		secretKey:     secretKey,
		expirationHrs: expirationHrs,
	}
}

func (tm *TokenManager) GenerarToken(postulanteID int, cedula string) (string, error) {
	claims := Claims{
		PostulanteID: postulanteID,
		Cedula:       cedula,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tm.expirationHrs) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", fmt.Errorf("error generando token: %w", err)
	}

	return tokenString, nil
}

func (tm *TokenManager) ValidarToken(tokenString string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(tm.secretKey), nil
	})

	if err != nil {
		return 0, "", domain.ErrTokenInvalido
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.PostulanteID, claims.Cedula, nil
	}

	return 0, "", domain.ErrTokenInvalido
}
