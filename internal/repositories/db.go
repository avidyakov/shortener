package repositories

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Link struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	OriginURL   string `gorm:"uniqueIndex"`
	ShortLinkID string `gorm:"uniqueIndex"`
	UserID      uint
}

type User struct {
	gorm.Model
	ID    uint `gorm:"primaryKey"`
	Links []Link
}

type DBRepo struct {
	db *gorm.DB
}

func (r *DBRepo) GetUrlsByUserId(userID int) ([]map[string]string, error) {
	var links []Link
	r.db.Find(&links, "user_id = ?", userID)
	var result []map[string]string
	for _, link := range links {
		result = append(result, map[string]string{
			"short_url":  link.ShortLinkID,
			"origin_url": link.OriginURL,
		})
	}
	return result, nil
}

func (r *DBRepo) CreateUser() (int, error) {
	user := &User{}
	tx := r.db.Create(user)
	return int(user.ID), tx.Error
}

func NewDBRepo(databaseDSN string) LinkRepo {
	sqlDB, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Link{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

	return &DBRepo{
		db: db,
	}
}

func (r *DBRepo) GetOriginLink(shortLinkID string) (originLink string, ok bool) {
	var link Link
	r.db.First(&link, "short_link_id = ?", shortLinkID)
	return link.OriginURL, link.ID == 0
}

func (r *DBRepo) GetShortLink(originLink string) (string, bool) {
	var link Link
	r.db.First(&link, "origin_url = ?", originLink)
	return link.ShortLinkID, link.ID == 0
}

func (r *DBRepo) CreateLink(shortLinkID string, originLink string, userID int) error {
	tx := r.db.Create(&Link{
		OriginURL:   originLink,
		ShortLinkID: shortLinkID,
		UserID:      uint(userID),
	})
	return tx.Error
}

func (r *DBRepo) RemoveLink(shortLinkID string) {
	r.db.Delete(&Link{}, "short_link_id = ?", shortLinkID)
}

func (r *DBRepo) CheckConnection() error {
	return r.db.Exec("SELECT 1").Error
}
