package repositories

import (
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

var repo LinkRepo

func TestMain(m *testing.M) {
	databaseDSN := "postgres://postgres:changeme@localhost:5432"
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating database...")
	db.Exec("CREATE DATABASE shortener_test")
	defer db.Exec("DROP DATABASE shortener_test")
	defer log.Println("Dropping database...")
	repo = NewDBRepo(databaseDSN)
	m.Run()
}

func TestCreateLink(t *testing.T) {
	repo.CreateLink("12345678", "https://www.google.com")
	originLink, ok := repo.GetLink("12345678")

	require.True(t, ok, "Expected to get link")
	require.Equal(t, "https://www.google.com", originLink)
}

func TestRemoveLink(t *testing.T) {
	repo.RemoveLink("12345678")
	_, ok := repo.GetLink("12345678")

	require.False(t, ok, "Expected to not get link")
}
