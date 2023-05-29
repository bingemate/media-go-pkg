package repository

import (
	"time"
)

type (
	MediaType       string
	WatchListStatus string
)

const (
	MediaTypeMovie             MediaType       = "Movie"
	MediaTypeTvShow            MediaType       = "TvShow"
	MediaTypeEpisode           MediaType       = "Episode"
	WatchListStatusPlanToWatch WatchListStatus = "PLAN_TO_WATCH"
	WatchListStatusWatching    WatchListStatus = "WATCHING"
	WatchListStatusFinished    WatchListStatus = "FINISHED"
	WatchListStatusAbandoned   WatchListStatus = "ABANDONED"
)

type Model struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MediaFile struct {
	Model
	Filename  string
	Duration  float64
	Audio     []Audio    `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
	Subtitles []Subtitle `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Media struct {
	ID          int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	MediaType   MediaType `gorm:"index"`
	ReleaseDate time.Time `gorm:"type:date"`
	Name        string
	TvShows     []TvShow   `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Movies      []Movie    `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Episodes    []Episode  `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
	Categories  []Category `gorm:"many2many:category_media;constraint:OnDelete:CASCADE;"`
}

type TvShow struct {
	Model
	Name     string
	MediaID  int       `gorm:"not null"`
	Media    Media     `gorm:"reference:MediaID"`
	Episodes []Episode `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
}

type Episode struct {
	Model
	Name        string
	NbEpisode   int
	NbSeason    int
	MediaID     int       `gorm:"not null"`
	Media       Media     `gorm:"reference:MediaID"`
	TvShowID    string    `gorm:"type:uuid;not null"`
	TvShow      TvShow    `gorm:"reference:TvShowID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Movie struct {
	Model
	Name        string
	MediaID     int       `gorm:"not null"`
	Media       Media     `gorm:"reference:MediaID"`
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID;constraint:OnDelete:CASCADE;"`
}

type Audio struct {
	Model
	Filename    string
	Language    string
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Subtitle struct {
	Model
	Filename    string
	Language    string
	MediaFileID string    `gorm:"type:uuid;not null"`
	MediaFile   MediaFile `gorm:"reference:MediaFileID"`
}

type Category struct {
	Model
	Name   string  `gorm:"uniqueIndex"`
	Medias []Media `gorm:"many2many:category_media"`
}

type CategoryMedia struct {
	MediaID    int      `gorm:"primaryKey"`
	Media      Media    `gorm:"reference:MediaID;constraint:OnDelete:CASCADE;"`
	CategoryID string   `gorm:"type:uuid;primaryKey"`
	Category   Category `gorm:"reference:CategoryID;constraint:OnDelete:CASCADE;"`
}

type Rating struct {
	UserID    string    `gorm:"type:uuid;primaryKey"`
	MediaID   int       `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Media     Media     `gorm:"reference:MediaID;constraint:OnDelete:CASCADE;"`
	Rating    int
}

type WatchListItems struct {
	UserID  string          `gorm:"type:uuid;primaryKey"`
	MediaID int             `gorm:"primaryKey"`
	Status  WatchListStatus `gorm:"index"`
}

func (WatchListItems) TableName() string {
	return "watch_list_item"
}
