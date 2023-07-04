package repositories

import (
	"github.com/avidyakov/shortener/internal/config"
	"testing"
)

var repo LinkRepo

func TestMain(m *testing.M) {
	config.Conf = config.NewConfig()
	repo = NewFileRepo()
	m.Run()
}

func TestCreateLink(t *testing.T) {
	repo.CreateLink("12345678", "https://www.google.com")
	originLink, ok := repo.GetLink("12345678")
	if !ok {
		t.Error("Expected to get link")
	}
	if originLink != "https://www.google.com" {
		t.Errorf("Expected to get %s, but got %s", "https://www.google.com", originLink)
	}
}

func TestRemoveLink(t *testing.T) {
	repo.RemoveLink("12345678")
	_, ok := repo.GetLink("12345678")
	if ok {
		t.Error("Expected to not get link")
	}
}
