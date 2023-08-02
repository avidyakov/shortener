package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Response struct {
	ShortURL  string `json:"short_url"`
	OriginURL string `json:"origin_url"`
}

func (h *Handlers) UserURLs(w http.ResponseWriter, r *http.Request) {

}

func (h *Handlers) getUserId(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(h.secretKey), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token os valid")
	return claims.UserID
}
