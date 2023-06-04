package handlers

import (
	"github.com/avidyakov/shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRedirect(t *testing.T) {
	repo = repositories.NewMemoryLink()

	repo.CreateLink("12345678", "https://www.google.com")
	req := httptest.NewRequest(http.MethodGet, "/12345678", nil)
	res := httptest.NewRecorder()

	redirect(res, req)

	require.Equal(t, http.StatusTemporaryRedirect, res.Code, "Expected status code %d, but got %d", http.StatusTemporaryRedirect, res.Code)
	require.Equal(t, "https://www.google.com", res.Header().Get("Location"), "Expected Location header %s, but got %s", "https://www.google.com", res.Header().Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
	repo = repositories.NewMemoryLink()

	req := httptest.NewRequest(http.MethodGet, "/12345678", nil)
	res := httptest.NewRecorder()

	redirect(res, req)

	require.Equal(t, http.StatusNotFound, res.Code, "Expected status code %d, but got %d", http.StatusNotFound, res.Code)
}

func TestCreateLink(t *testing.T) {
	repo = repositories.NewMemoryLink()

	req := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader("https://www.google.com"),
	)
	res := httptest.NewRecorder()

	createShortLink(res, req)

	require.Equal(t, http.StatusCreated, res.Code, "Expected status code %d, but got %d", http.StatusCreated, res.Code)
}
