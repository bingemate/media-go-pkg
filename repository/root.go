package repository

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&MediaFile{},
		&TvShow{},
		&Episode{},
		&Movie{},
		&Audio{},
		&Subtitle{},
		&Category{},
		&CategoryMovie{},
		&CategoryTvShow{},
		&MovieRating{},
		&TvShowRating{},
		&MovieComment{},
		&TvShowComment{},
		&MovieWatchListItem{},
		&TvShowWatchListItem{},
	)
}
