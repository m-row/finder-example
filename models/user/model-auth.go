package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// jwt structs

type Token struct {
	Value     string `json:"value"`
	Type      string `json:"type"`
	ExpiresIn string `json:"expires_in"`
}

type TokenResponse struct {
	Token Token `json:"token"`
}

type FileTokenResponse struct {
	Token Token  `json:"token"`
	URL   string `json:"url"`
}

type CustomOrderFileClaims struct {
	OrderID string `json:"order_id"`
	jwt.RegisteredClaims
}

// jwt stuff

func (m *Model) ParseToken(
	bearer *string,
) (*jwt.Token, *jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		*bearer,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return "secret", nil
		},
	)
	if err != nil {
		return &jwt.Token{}, &jwt.RegisteredClaims{}, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return token, claims, nil
	}
	return &jwt.Token{}, &jwt.RegisteredClaims{}, err
}

func (m *Model) GenTokenResponse() (TokenResponse, error) {
	// Create the claims
	claims := jwt.RegisteredClaims{
		Issuer:   "localhost",
		Subject:  m.ID.String(),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(5 * time.Hour),
		),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString("secret")
	if err != nil {
		return TokenResponse{}, err
	}
	tokenResponse := TokenResponse{
		Token: Token{
			Value:     signedString,
			Type:      "bearer",
			ExpiresIn: fmt.Sprintf("%dh", int(5)),
		},
	}
	return tokenResponse, nil
}

// GenCookie generates an http only cookie for the token given
// with expires time in the future for valid tokens and
// in the past for invalidating tokens (logging out).
func (m *Model) GenCookie(token Token, expires time.Time) http.Cookie {
	return http.Cookie{
		Name:     "accessToken",
		Value:    token.Value,
		Path:     "/",
		MaxAge:   0,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
