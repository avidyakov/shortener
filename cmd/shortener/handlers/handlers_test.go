package handlers

import (
	"github.com/avidyakov/shortener/cmd/shortener/config"
	"github.com/avidyakov/shortener/cmd/shortener/repositories"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	config.Cfg = &config.Config{
		BaseURL:    "http://localhost:8080",
		ServerAddr: ":8080",
	}
}

func TestRedirect(t *testing.T) {
	repo = repositories.NewMemoryLink()

	ts := httptest.NewServer(LinkRouter())
	defer ts.Close()
	repo.CreateLink("12345678", "https://www.google.com")
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Get(ts.URL + "/12345678")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Expected status code %d, but got %d", http.StatusTemporaryRedirect, resp.StatusCode)
	require.Equal(t, "https://www.google.com", resp.Header.Get("Location"), "Expected Location header %s, but got %s", "https://www.google.com", resp.Header.Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
	repo = repositories.NewMemoryLink()

	ts := httptest.NewServer(LinkRouter())
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Get(ts.URL + "/12345678")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
}

func TestCreateLink(t *testing.T) {
	repo = repositories.NewMemoryLink()

	ts := httptest.NewServer(LinkRouter())
	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Post(ts.URL, "text/plain", strings.NewReader("https://www.google.com"))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
}
