package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type Response struct {
	ShortURL  string `json:"short_url"`
	OriginURL string `json:"origin_url"`
}

func (h *Handlers) UserURLs(w http.ResponseWriter, r *http.Request) {
	parsedToken, _ := r.Cookie("token")
	userID := h.getUserID(parsedToken.Value)
	if userID == -1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	urls, err := h.repo.GetUrlsByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := []Response{}
	for i := range urls {
		response = append(response, Response{ShortURL: urls[i]["short_url"], OriginURL: urls[i]["origin_url"]})
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) getUserID(tokenString string) int {
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
