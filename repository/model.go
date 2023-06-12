package repository

import (
	"time"
)

type (
	WatchListStatus string
)

const (
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
	Size      int64
	Audios    []Audio    `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
	Subtitles []Subtitle `gorm:"foreignKey:MediaFileID;constraint:OnDelete:CASCADE;"`
}

//type Media struct {
//	ID          int       `gorm:"primaryKey"`
//	CreatedAt   time.Time `gorm:"autoCreateTime"`
//	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
//	MediaType   MediaType `gorm:"index"`
//	ReleaseDate time.Time `gorm:"type:date"`
//	Name        string
//	TvShows     []TvShow   `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
//	Movies      []Movie    `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
//	Episodes    []Episode  `gorm:"foreignKey:MediaID;constraint:OnDelete:CASCADE;"`
//	Categories  []Category `gorm:"many2many:category_media;constraint:OnDelete:CASCADE;"`
//}

type TvShow struct {
	ID          int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Name        string
	ReleaseDate time.Time       `gorm:"type:date"`
	Episodes    []Episode       `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
	Categories  []Category      `gorm:"many2many:category_tv_show;constraint:OnDelete:CASCADE;"`
	Ratings     []TvShowRating  `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
	Comments    []TvShowComment `gorm:"foreignKey:TvShowID;constraint:OnDelete:CASCADE;"`
}

type Episode struct {
	ID          int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Name        string
	NbEpisode   int
	NbSeason    int
	ReleaseDate time.Time  `gorm:"type:date"`
	TvShowID    int        `gorm:"not null"`
	TvShow      TvShow     `gorm:"reference:TvShowID"`
	MediaFileID *string    `gorm:"type:uuid"`
	MediaFile   *MediaFile `gorm:"reference:MediaFileID;constraint:OnDelete:SET NULL;"`
}

type Movie struct {
	ID          int       `gorm:"primaryKey"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Name        string
	ReleaseDate time.Time      `gorm:"type:date"`
	MediaFileID *string        `gorm:"type:uuid"`
	MediaFile   *MediaFile     `gorm:"reference:MediaFileID;constraint:OnDelete:SET NULL;"`
	Categories  []Category     `gorm:"many2many:category_movie;constraint:OnDelete:CASCADE;"`
	Ratings     []MovieRating  `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE;"`
	Comments    []MovieComment `gorm:"foreignKey:MovieID;constraint:OnDelete:CASCADE;"`
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
	Name   string   `gorm:"uniqueIndex"`
	Movies []Movie  `gorm:"many2many:category_movie;constraint:OnDelete:CASCADE;"`
	TvShow []TvShow `gorm:"many2many:category_tv_show;constraint:OnDelete:CASCADE;"`
}

type CategoryMovie struct {
	MovieID    int      `gorm:"primaryKey"`
	Movie      Movie    `gorm:"reference:MovieID;constraint:OnDelete:CASCADE;"`
	CategoryID string   `gorm:"type:uuid;primaryKey"`
	Category   Category `gorm:"reference:CategoryID;constraint:OnDelete:CASCADE;"`
}

func (CategoryMovie) TableName() string {
	return "category_movie"
}

type CategoryTvShow struct {
	TvShowID   int      `gorm:"primaryKey"`
	TvShow     TvShow   `gorm:"reference:TvShowID;constraint:OnDelete:CASCADE;"`
	CategoryID string   `gorm:"type:uuid;primaryKey"`
	Category   Category `gorm:"reference:CategoryID;constraint:OnDelete:CASCADE;"`
}

func (CategoryTvShow) TableName() string {
	return "category_tv_show"
}

//type Rating struct {
//	UserID    string    `gorm:"type:uuid;primaryKey"`
//	MediaID   int       `gorm:"primaryKey"`
//	CreatedAt time.Time `gorm:"autoCreateTime"`
//	UpdatedAt time.Time `gorm:"autoUpdateTime"`
//	Rating    int
//}

type MovieRating struct {
	UserID    string    `gorm:"type:uuid;primaryKey"`
	MovieID   int       `gorm:"primaryKey"`
	Movie     Movie     `gorm:"reference:MovieID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Rating    int
}

type TvShowRating struct {
	UserID    string    `gorm:"type:uuid;primaryKey"`
	TvShowID  int       `gorm:"primaryKey"`
	TvShow    TvShow    `gorm:"reference:TvShowID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Rating    int
}

//type Comment struct {
//	Model
//	UserID  string `gorm:"type:uuid;not null"`
//	MediaID int    `gorm:"not null"`
//	Content string
//}

type MovieComment struct {
	Model
	Content string
	UserID  string `gorm:"type:uuid;not null"`
	MovieID int    `gorm:"not null"`
	Movie   Movie  `gorm:"reference:MovieID;constraint:OnDelete:CASCADE;"`
}

type TvShowComment struct {
	Model
	Content  string
	UserID   string `gorm:"type:uuid;not null"`
	TvShowID int    `gorm:"not null"`
	TvShow   TvShow `gorm:"reference:TvShowID;constraint:OnDelete:CASCADE;"`
}

//type WatchListItem struct {
//	UserID  string          `gorm:"type:uuid;primaryKey"`
//	MediaID int             `gorm:"primaryKey"`
//	Status  WatchListStatus `gorm:"index"`
//}
//
//func (WatchListItem) TableName() string {
//	return "watch_list_item"
//}

type MovieWatchListItem struct {
	UserID  string          `gorm:"type:uuid;primaryKey"`
	MovieID int             `gorm:"primaryKey"`
	Movie   Movie           `gorm:"reference:MovieID;constraint:OnDelete:CASCADE;"`
	Status  WatchListStatus `gorm:"index;not null"`
}

func (MovieWatchListItem) TableName() string {
	return "movie_watch_list_item"
}

type TvShowWatchListItem struct {
	UserID   string          `gorm:"type:uuid;primaryKey"`
	TvShowID int             `gorm:"primaryKey"`
	TvShow   TvShow          `gorm:"reference:TvShowID;constraint:OnDelete:CASCADE;"`
	Status   WatchListStatus `gorm:"index;not null"`
}

func (TvShowWatchListItem) TableName() string {
	return "tv_show_watch_list_item"
}
