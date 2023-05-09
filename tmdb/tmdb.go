package tmdb

import "github.com/ryanbradynd05/go-tmdb"

const imageBaseURL = "https://image.tmdb.org/t/p/original"
const emptyProfileURL = "https://bingemate.fr/assets/empty_profile.jpg"
const emptyBackdropURL = "https://bingemate.fr/assets/empty_background.jpg"
const emptyPosterURL = "https://bingemate.fr/assets/empty_poster.jpg"

// Genre represents a movie/TV genre with its ID and name.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Person represents a person involved in a movie or TV show with their ID, character name,
// real name, and the URL to their profile picture.
type Person struct {
	ID         int    `json:"id"`
	Character  string `json:"character"`
	Name       string `json:"name"`
	ProfileURL string `json:"profile_url"`
}

// Studio represents a movie/TV studio with its ID, name, and logo URL.
type Studio struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logo_url"`
}

// Movie represents a movie with its attributes such as ID, actors list (Person), backdrop URL,
// crew list (Person), genre list (Genre), overview, poster URL, release date, studio list (Studio),
// title, vote average, and vote count.
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

// TVEpisode represents a TV episode with its attributes such as ID, TV show ID, poster URL,
// season number, episode number, name, overview, and air date.
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

// TVShow represents a TV show with its attributes such as ID, actors list (Person), backdrop URL,
// crew list (Person), genre list (Genre), overview, poster URL, release date, studio list (Studio),
// status, next episode (TVEpisode), title, seasons count, vote average, and vote count.
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

// MediaClient is an interface for a media client API.
type MediaClient interface {
	GetMovie(id int) (*Movie, error)
	GetTVShow(id int) (*TVShow, error)
	GetTVEpisode(tvId, season, episodeNumber int) (*TVEpisode, error)
	GetTVSeasonEpisodes(id int, season int) (*[]TVEpisode, error)
	GetTrendingMovies() (*[]Movie, error)
	GetTrendingTVShows() (*[]TVShow, error)
	GetRecentMovies() (*[]Movie, error)
	GetRecentTVShows() (*[]TVShow, error)
	SearchMovies(query string) (*[]Movie, error)
	SearchTVShows(query string) (*[]TVShow, error)
	GetMoviesByGenre(genreID int) (*[]Movie, error)
	GetTVShowsByGenre(genreID int) (*[]TVShow, error)
	GetMoviesByActor(actorID int) (*[]Movie, error)
	GetTVShowsByActor(actorID int) (*[]TVShow, error)
	GetMoviesByDirector(directorID int) (*[]Movie, error)
	GetTVShowsByDirector(directorID int) (*[]TVShow, error)
	GetMoviesByStudio(studioID int) (*[]Movie, error)
	GetTVShowsByStudio(studioID int) (*[]TVShow, error)
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
		BackdropURL: backdropImgURL(movie.BackdropPath),
		Crew:        *extractMovieCrew(credits),
		Genres:      *extractGenres(&movie.Genres),
		Overview:    movie.Overview,
		PosterURL:   posterImgURL(movie.PosterPath),
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
		BackdropURL: backdropImgURL(tvShow.BackdropPath),
		Crew:        *extractTVCrew(credits),
		Genres:      *extractGenres(&tvShow.Genres),
		Overview:    tvShow.Overview,
		PosterURL:   posterImgURL(tvShow.PosterPath),
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
				PosterURL:     posterImgURL(tvShow.NextEpisodeToAir.StillPath),
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
		PosterURL:     posterImgURL(episode.StillPath),
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
			PosterURL:     posterImgURL(episode.StillPath),
			EpisodeNumber: episode.EpisodeNumber,
			SeasonNumber:  episode.SeasonNumber,
			Name:          episode.Name,
			Overview:      episode.Overview,
			AirDate:       episode.AirDate,
		}
	}
	return &extractedEpisodes, nil
}

func (m *mediaClient) GetTrendingMovies() (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetTrendingTVShows() (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetRecentMovies() (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetRecentTVShows() (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) SearchMovies(query string) (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) SearchTVShows(query string) (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetMoviesByGenre(genreID int) (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetTVShowsByGenre(genreID int) (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetMoviesByActor(actorID int) (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetTVShowsByActor(actorID int) (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetMoviesByDirector(directorID int) (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetTVShowsByDirector(directorID int) (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetMoviesByStudio(studioID int) (*[]Movie, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetTVShowsByStudio(studioID int) (*[]TVShow, error) {
	//TODO implement me
	panic("implement me")
}

// extractMovieActors extracts actors from movie credits and returns a list of Person.
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

// extractTVActors extracts actors from TV show credits and returns a list of Person.
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

// extractMovieCrew extracts crew from movie credits and returns a list of Person.
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

// extractTVCrew extracts crew from TV show credits and returns a list of Person.
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

// extractGenres extracts genres from a list of genre structs and returns a list of Genre.
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

// extractStudios extracts studios from a list of studio structs and returns a list of Studio.
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

// profileImgURL returns the URL of a profile image given its path or an empty string if the path is empty.
func profileImgURL(path string) string {
	if path == "" {
		return emptyProfileURL
	}
	return imageBaseURL + path
}

// backdropImgURL returns the URL of a backdrop image given its path or an empty string if the path is empty.
func backdropImgURL(path string) string {
	if path == "" {
		return emptyBackdropURL
	}
	return imageBaseURL + path
}

// posterImgURL returns the URL of a poster image given its path or an empty string if the path is empty.
func posterImgURL(path string) string {
	if path == "" {
		return emptyPosterURL
	}
	return imageBaseURL + path
}
