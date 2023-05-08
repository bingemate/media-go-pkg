package tmdb

import "github.com/ryanbradynd05/go-tmdb"

const imageBaseURL = "https://image.tmdb.org/t/p/original"
const blankProfileURL = "https://bingemate.fr/assets/images/blank-profile.png"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Person struct {
	ID         int    `json:"id"`
	Character  string `json:"character"`
	Name       string `json:"name"`
	ProfileURL string `json:"profile_url"`
}

type Studio struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}

type Movie struct {
	ID          int      `json:"id"`
	Actors      []Person `json:"actors"`
	BackdropURL string   `json:"backdrop_url"`
	Crew        []Person `json:"crew"`
	Genres      []Genre  `json:"genres"`
	Overview    string   `json:"overview"`
	PosterURL   string   `json:"poster_url"`
	ReleaseDate string   `json:"release_date"`
	Studios     []Studio `json:"studios"`
	Title       string   `json:"title"`
	VoteAverage float32  `json:"vote_average"`
	VoteCount   int      `json:"vote_count"`
}

type TVEpisode struct {
	ID            int    `json:"id"`
	TVShowID      int    `json:"tv_show_id"`
	PosterURL     string `json:"poster_url"`
	EpisodeNumber int    `json:"episode_number"`
	SeasonNumber  int    `json:"season_number"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	AirDate       string `json:"air_date"`
}

type TVShow struct {
	ID           int        `json:"id"`
	Actors       []Person   `json:"actors"`
	BackdropURL  string     `json:"backdrop_url"`
	Crew         []Person   `json:"crew"`
	Genres       []Genre    `json:"genres"`
	Overview     string     `json:"overview"`
	PosterURL    string     `json:"poster_url"`
	ReleaseDate  string     `json:"release_date"`
	Studios      []Studio   `json:"studios"`
	Status       string     `json:"status"`
	NextEpisode  *TVEpisode `json:"next_episode"`
	Title        string     `json:"title"`
	SeasonsCount int        `json:"seasons_count"`
	VoteAverage  float32    `json:"vote_average"`
	VoteCount    int        `json:"vote_count"`
}

type MediaClient interface {
	GetMovie(id int) (*Movie, error)
	GetTVShow(id int) (*TVShow, error)
	GetTVEpisode(tvId, season, episodeNumber int) (*TVEpisode, error)
	GetTVSeasonEpisodes(id int, season int) (*[]TVEpisode, error)
}

type mediaClient struct {
	tmdbClient *tmdb.TMDb
	options    map[string]string
}

func NewMediaClient(apiKey string) MediaClient {
	config := tmdb.Config{
		APIKey:   apiKey,
		Proxies:  nil,
		UseProxy: false,
	}
	return &mediaClient{
		tmdbClient: tmdb.Init(config),
		options: map[string]string{
			"language": "fr",
		},
	}
}

func (m *mediaClient) GetMovie(id int) (*Movie, error) {
	movie, err := m.tmdbClient.GetMovieInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	credits, err := m.tmdbClient.GetMovieCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	return &Movie{
		ID:          movie.ID,
		Actors:      *extractMovieActors(credits),
		BackdropURL: imageBaseURL + movie.BackdropPath,
		Crew:        *extractMovieCrew(credits),
		Genres:      *extractGenres(&movie.Genres),
		Overview:    movie.Overview,
		PosterURL:   imageBaseURL + movie.PosterPath,
		ReleaseDate: movie.ReleaseDate,
		Studios:     *extractStudios(&movie.ProductionCompanies),
		Title:       movie.Title,
		VoteAverage: movie.VoteAverage,
		VoteCount:   int(movie.VoteCount),
	}, err
}

func (m *mediaClient) GetTVShow(id int) (*TVShow, error) {
	tvShow, err := m.tmdbClient.GetTvInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	credits, err := m.tmdbClient.GetTvCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	return &TVShow{
		ID:          tvShow.ID,
		Actors:      *extractTVActors(credits),
		BackdropURL: imageBaseURL + tvShow.BackdropPath,
		Crew:        *extractTVCrew(credits),
		Genres:      *extractGenres(&tvShow.Genres),
		Overview:    tvShow.Overview,
		PosterURL:   imageBaseURL + tvShow.PosterPath,
		ReleaseDate: tvShow.FirstAirDate,
		Studios:     *extractStudios(&tvShow.ProductionCompanies),
		Status:      tvShow.Status,
		Title:       tvShow.Name,
		NextEpisode: func() *TVEpisode {
			if tvShow.NextEpisodeToAir.ID == 0 {
				return nil
			}
			return &TVEpisode{
				ID:            tvShow.NextEpisodeToAir.ID,
				TVShowID:      tvShow.ID,
				PosterURL:     imageBaseURL + tvShow.NextEpisodeToAir.StillPath,
				EpisodeNumber: tvShow.NextEpisodeToAir.EpisodeNumber,
				SeasonNumber:  tvShow.NextEpisodeToAir.SeasonNumber,
				Name:          tvShow.NextEpisodeToAir.Name,
				Overview:      tvShow.NextEpisodeToAir.Overview,
				AirDate:       tvShow.NextEpisodeToAir.AirDate,
			}
		}(),
		SeasonsCount: tvShow.NumberOfSeasons,
		VoteAverage:  tvShow.VoteAverage,
		VoteCount:    int(tvShow.VoteCount),
	}, nil
}

func (m *mediaClient) GetTVEpisode(tvId, season, episodeNumber int) (*TVEpisode, error) {
	episode, err := m.tmdbClient.GetTvEpisodeInfo(tvId, season, episodeNumber, m.options)
	if err != nil {
		return nil, err
	}
	return &TVEpisode{
		ID:            episode.ID,
		TVShowID:      tvId,
		PosterURL:     imageBaseURL + episode.StillPath,
		EpisodeNumber: episode.EpisodeNumber,
		SeasonNumber:  episode.SeasonNumber,
		Name:          episode.Name,
		Overview:      episode.Overview,
		AirDate:       episode.AirDate,
	}, nil
}

func (m *mediaClient) GetTVSeasonEpisodes(tvId int, season int) (*[]TVEpisode, error) {
	episodes, err := m.tmdbClient.GetTvSeasonInfo(tvId, season, m.options)
	if err != nil {
		return nil, err
	}
	var extractedEpisodes = make([]TVEpisode, len(episodes.Episodes))
	for i, episode := range episodes.Episodes {
		extractedEpisodes[i] = TVEpisode{
			ID:            episode.ID,
			TVShowID:      tvId,
			EpisodeNumber: episode.EpisodeNumber,
			SeasonNumber:  episode.SeasonNumber,
			Name:          episode.Name,
			Overview:      episode.Overview,
			AirDate:       episode.AirDate,
		}
	}
	return &extractedEpisodes, nil
}

func extractMovieActors(credits *tmdb.MovieCredits) *[]Person {
	var actors = make([]Person, len(credits.Cast))
	for i, cast := range credits.Cast {
		actors[i] = Person{
			ID:         cast.ID,
			Character:  cast.Character,
			Name:       cast.Name,
			ProfileURL: profileImgURL(cast.ProfilePath),
		}
	}
	return &actors
}

func extractTVActors(credits *tmdb.TvCredits) *[]Person {
	var actors = make([]Person, len(credits.Cast))
	for i, cast := range credits.Cast {
		actors[i] = Person{
			ID:         cast.ID,
			Character:  cast.Character,
			Name:       cast.Name,
			ProfileURL: profileImgURL(cast.ProfilePath),
		}
	}
	return &actors
}

func extractMovieCrew(credits *tmdb.MovieCredits) *[]Person {
	var crew = make([]Person, len(credits.Crew))
	for i, cast := range credits.Crew {
		crew[i] = Person{
			ID:         cast.ID,
			Character:  cast.Job,
			Name:       cast.Name,
			ProfileURL: profileImgURL(cast.ProfilePath),
		}
	}
	return &crew
}

func extractTVCrew(credits *tmdb.TvCredits) *[]Person {
	var crew = make([]Person, len(credits.Crew))
	for i, cast := range credits.Crew {
		crew[i] = Person{
			ID:         cast.ID,
			Character:  cast.Job,
			Name:       cast.Name,
			ProfileURL: profileImgURL(cast.ProfilePath),
		}
	}
	return &crew
}

func extractGenres(genres *[]struct {
	ID   int
	Name string
}) *[]Genre {
	var extractedGenres = make([]Genre, len(*genres))
	for i, genre := range *genres {
		extractedGenres[i] = Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}
	}
	return &extractedGenres
}

func extractStudios(studios *[]struct {
	ID        int
	Name      string
	LogoPath  string `json:"logo_path"`
	Iso3166_1 string `json:"origin_country"`
}) *[]Studio {
	var extractedStudios = make([]Studio, len(*studios))
	for i, studio := range *studios {
		extractedStudios[i] = Studio{
			ID:      studio.ID,
			Name:    studio.Name,
			LogoURL: profileImgURL(studio.LogoPath),
		}
	}
	return &extractedStudios
}

func profileImgURL(path string) string {
	if path == "" {
		return blankProfileURL
	}
	return imageBaseURL + path
}
