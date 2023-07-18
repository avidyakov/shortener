package repositories

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Link struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	OriginURL   string `gorm:"uniqueIndex"`
	ShortLinkID string `gorm:"uniqueIndex"`
}

type dbRepo struct {
	db *gorm.DB
}

func NewDBRepo(databaseDSN string) LinkRepo {
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Link{})
	if err != nil {
		log.Fatal(err)
	}

	return &dbRepo{
		db: db,
	}
}

func (r *dbRepo) GetOriginLink(shortLinkID string) (originLink string, ok bool) {
	var link Link
	r.db.First(&link, "short_link_id = ?", shortLinkID)

	if link.ID == 0 {
		return "", false
	}
	return link.OriginURL, true
}

func (r *dbRepo) GetShortLink(originLink string) (string, bool) {
	var link Link
	r.db.First(&link, "origin_url = ?", originLink)

	if link.ID == 0 {
		return "", false
	}
	return link.ShortLinkID, true
}

func (r *dbRepo) CreateLink(shortLinkID string, originLink string) error {
	tx := r.db.Create(&Link{
		OriginURL:   originLink,
		ShortLinkID: shortLinkID,
	})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *dbRepo) RemoveLink(shortLinkID string) {
	r.db.Delete(&Link{}, "short_link_id = ?", shortLinkID)
}
