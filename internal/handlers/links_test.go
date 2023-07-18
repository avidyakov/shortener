package handlers

import (
	"encoding/json"
	"github.com/avidyakov/shortener/internal/logger"
	"github.com/avidyakov/shortener/internal/models"
	"github.com/avidyakov/shortener/internal/repositories"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var h *Handlers
var ts *httptest.Server
var client *http.Client
var repo repositories.LinkRepo

func TestMain(m *testing.M) {
	logger.Log, _ = zap.NewProduction()
	defer logger.Log.Sync()

	repo = repositories.NewMemoryRepo()
	h = NewHandlers(repo, "http://localhost:8080", "")

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

func TestCreateLinkJSON(t *testing.T) {
	resp, err := client.Post(ts.URL, "application/json", strings.NewReader(`{"url": "https://www.google.com"}`))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
}

func TestCreateLinkInvalidURL(t *testing.T) {
	resp, err := client.Post(ts.URL, "text/plain", strings.NewReader("invalid url"))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code %d, but got %d", http.StatusBadRequest, resp.StatusCode)
}

func TestCreateLinkBacket(t *testing.T) {
	resp, err := client.Post(ts.URL+"/api/shorten/batch", "application/json", strings.NewReader(`[{"original_url": "https://www.google.com", "correlation_id": "1"}, {"original_url": "https://ya.ru", "correlation_id": "2"}]`))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	links := make([]models.ResponseLinkBatch, 2)
	err = json.NewDecoder(resp.Body).Decode(&links)
	require.NoError(t, err)
	require.Len(t, links, 2)

	shortGoogleID := strings.Replace(links[0].ShortURL, h.baseURL+"/", "", -1)
	google, _ := repo.GetLink(shortGoogleID)
	require.Equal(t, google, "https://www.google.com")

	shortYandexID := strings.Replace(links[1].ShortURL, h.baseURL+"/", "", -1)
	yandex, _ := repo.GetLink(shortYandexID)
	require.Equal(t, yandex, "https://ya.ru")
}
