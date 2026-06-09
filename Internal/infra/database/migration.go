package database

import (
	"errors"
	"log"

	urlEntity "github.com/haramurti/valeth-urls-shortener/Internal/app/url/entity"
	userEntity "github.com/haramurti/valeth-urls-shortener/Internal/app/user/entity"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&userEntity.User{},
		&urlEntity.Url{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database")
		return errors.New("Failed to smigrate database")
	}
	log.Println("Database Migrated Successfully")

	return nil
}

//this is acomment
