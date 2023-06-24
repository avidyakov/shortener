package handlers

import (
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var h *LinkHandlers
var ts *httptest.Server
var client *http.Client

func TestMain(m *testing.M) {
	logger.Log, _ = zap.NewProduction()
	defer logger.Log.Sync()

	h = NewLinkHandlers(
		repositories.NewMemoryLink(), "http://localhost:8080",
	)

	ts = httptest.NewServer(h.LinkRouter())
	defer ts.Close()

	client = ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	m.Run()
}

func createLink(t *testing.T, shortLink, originLink string) {
	h.repo.CreateLink(shortLink, originLink)

	t.Cleanup(func() {
		h.repo.RemoveLink(shortLink)
	})
}

func TestRedirect(t *testing.T) {
	createLink(t, "12345678", "https://www.google.com")

	resp, err := client.Get(ts.URL + "/12345678")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode, "Expected status code %d, but got %d", http.StatusTemporaryRedirect, resp.StatusCode)
	require.Equal(t, "https://www.google.com", resp.Header.Get("Location"), "Expected Location header %s, but got %s", "https://www.google.com", resp.Header.Get("Location"))
}

func TestRedirectNotFound(t *testing.T) {
	resp, err := client.Get(ts.URL + "/12345678")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
}

func TestCreateLink(t *testing.T) {
	resp, err := client.Post(ts.URL, "text/plain", strings.NewReader("https://www.google.com"))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
}
