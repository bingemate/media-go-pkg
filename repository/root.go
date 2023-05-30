package repository

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&MediaFile{},
		&Media{},
		&TvShow{},
		&Episode{},
		&Movie{},
		&Audio{},
		&Subtitle{},
		&Category{},
		&CategoryMedia{},
		&Rating{},
		&WatchListItem{},
		&Comment{},
	)
}
