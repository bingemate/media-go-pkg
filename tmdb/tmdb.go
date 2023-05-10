package tmdb

import (
	"github.com/ryanbradynd05/go-tmdb"
	"sort"
	"strconv"
	"time"
)

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

// TVEpisodeRelease represents a TV episode release with its attributes such as ID, name, episode
type TVEpisodeRelease struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	EpisodeNumber int    `json:"episode_number"`
	SeasonNumber  int    `json:"season_number"`
	TVShowName    string `json:"tv_show_name"`
	AirDate       string `json:"air_date"`
}

// MovieRelease represents a movie release with its attributes such as ID, title, and release date.
type MovieRelease struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
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
	Networks     []Studio   `json:"networks"`
	Status       string     `json:"status"`
	NextEpisode  *TVEpisode `json:"next_episode"`
	Title        string     `json:"title"`
	SeasonsCount int        `json:"seasons_count"`
	VoteAverage  float32    `json:"vote_average"`
	VoteCount    int        `json:"vote_count"`
}

//TODO Create a result struct wrapper that contains the results and the total pages

// MediaClient is an interface for a media client API.
type MediaClient interface {
	GetMovie(id int) (*Movie, error)
	GetTVShow(id int) (*TVShow, error)
	GetTVEpisode(tvId, season, episodeNumber int) (*TVEpisode, error)
	GetTVSeasonEpisodes(id int, season int) (*[]TVEpisode, error)
	GetPopularMovies(page int) (*[]Movie, error)
	GetPopularTVShows(page int) (*[]TVShow, error)
	GetRecentMovies() (*[]Movie, error)
	GetRecentTVShows() (*[]TVShow, error)
	SearchMovies(query string, page int) (*[]Movie, error)
	SearchTVShows(query string, page int) (*[]TVShow, error)
	GetMoviesByGenre(genreID int, page int) (*[]Movie, error)
	GetTVShowsByGenre(genreID int, page int) (*[]TVShow, error)
	GetMoviesByActor(actorID int, page int) (*[]Movie, error)
	GetTVShowsByActor(actorID int, page int) (*[]TVShow, error)
	GetMoviesByDirector(directorID int, page int) (*[]Movie, error)
	GetTVShowsByDirector(directorID int, page int) (*[]TVShow, error)
	GetMoviesByStudio(studioID int, page int) (*[]Movie, error)
	GetTVShowsByNetwork(studioID int, page int) (*[]TVShow, error)
	GetTVShowsNextEpisode(tvIds []int, startDate, endDate string) (*[]TVEpisodeRelease, error)
	GetMoviesReleases(movieIds []int, startDate, endDate string) (*[]MovieRelease, error)
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

// GetMovie retrieves movie info and credits by ID and returns a Movie object.
func (m *mediaClient) GetMovie(id int) (*Movie, error) {
	movie, err := m.tmdbClient.GetMovieInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	credits, err := m.tmdbClient.GetMovieCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	return extractMovie(movie, credits), nil
}

// GetTVShow retrieves TV show info and credits by ID and returns a TVShow object.
func (m *mediaClient) GetTVShow(id int) (*TVShow, error) {
	tvShow, err := m.tmdbClient.GetTvInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	credits, err := m.tmdbClient.GetTvCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	return extractTVShow(tvShow, credits), nil
}

// GetTVEpisode retrieves the information of a TV episode by TV show ID, season number and episode number and returns a TVEpisode object.
func (m *mediaClient) GetTVEpisode(tvId, season, episodeNumber int) (*TVEpisode, error) {
	episode, err := m.tmdbClient.GetTvEpisodeInfo(tvId, season, episodeNumber, m.options)
	if err != nil {
		return nil, err
	}
	return extractTVEpisode(tvId, episode), nil
}

// GetTVSeasonEpisodes retrieves all episodes from a TV show season and returns a slice of TVEpisode objects.
func (m *mediaClient) GetTVSeasonEpisodes(tvId int, season int) (*[]TVEpisode, error) {
	episodes, err := m.tmdbClient.GetTvSeasonInfo(tvId, season, m.options)
	if err != nil {
		return nil, err
	}
	var extractedEpisodes = make([]TVEpisode, len(episodes.Episodes))
	for i, episode := range episodes.Episodes {
		extractedEpisodes[i] = *extractTVEpisode(tvId, &episode)
	}
	return &extractedEpisodes, nil
}

func (m *mediaClient) GetPopularMovies(page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	movies, err := m.tmdbClient.GetMoviePopular(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) GetPopularTVShows(page int) (*[]TVShow, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	tvShows, err := m.tmdbClient.GetTvPopular(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowShort(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetRecentMovies() (*[]Movie, error) {
	options := extractOptions(m.options)
	options["region"] = "fr"
	movies := make([]tmdb.MovieShort, 0)
	// Get the 100 most recent movies in France (20 per page)
	for page := 1; page <= 5; page++ {
		options["page"] = strconv.Itoa(page)
		retrievedMovies, err := m.tmdbClient.GetMovieNowPlaying(options)
		if err != nil {
			return nil, err
		}
		movies = append(movies, retrievedMovies.Results...)
	}
	// Sort them by popularity
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Popularity > movies[j].Popularity
	})
	var extractedMovies = make([]Movie, 0)
	// Get the 20 most popular
	for _, movie := range movies[:20] {
		extractedMovies = append(extractedMovies, *extractMovieShort(&movie))
	}
	// Sort them by release date (the most recent first)
	sort.Slice(extractedMovies, func(i, j int) bool {
		releaseDateI, err := time.Parse("2006-01-02", extractedMovies[i].ReleaseDate)
		if err != nil {
			return false
		}
		releaseDateJ, err := time.Parse("2006-01-02", extractedMovies[j].ReleaseDate)
		if err != nil {
			return false
		}
		return releaseDateI.After(releaseDateJ)
	})
	// Return result
	return &extractedMovies, nil
}

func (m *mediaClient) GetRecentTVShows() (*[]TVShow, error) {
	options := extractOptions(m.options)
	tvshows := make([]tmdb.TvShort, 0)
	// Get the 100 most recent tvshows in France (20 per page)
	for page := 1; page <= 5; page++ {
		options["page"] = strconv.Itoa(page)
		retrievedTVShows, err := m.tmdbClient.GetTvAiringToday(options)
		if err != nil {
			return nil, err
		}
		tvshows = append(tvshows, retrievedTVShows.Results...)
	}
	// Sort them by popularity
	sort.Slice(tvshows, func(i, j int) bool {
		return tvshows[i].Popularity > tvshows[j].Popularity
	})
	var extractedTVShows = make([]TVShow, 0)
	// Get the 20 most popular
	for _, tvshow := range tvshows[:20] {
		extractedTVShows = append(extractedTVShows, *extractTVShowShort(&tvshow))
	}
	// Return result
	return &extractedTVShows, nil
}

func (m *mediaClient) SearchMovies(query string, page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["region"] = "fr"
	movies, err := m.tmdbClient.SearchMovie(query, options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) SearchTVShows(query string, page int) (*[]TVShow, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	tvShows, err := m.tmdbClient.SearchTv(query, options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowResult(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetMoviesByGenre(genreID int, page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_genres"] = strconv.Itoa(genreID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) GetTVShowsByGenre(genreID int, page int) (*[]TVShow, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_genres"] = strconv.Itoa(genreID)
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowShort(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetMoviesByActor(actorID int, page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_cast"] = strconv.Itoa(actorID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) GetTVShowsByActor(actorID int, page int) (*[]TVShow, error) {
	//TODO: Test this function, not sure "with_cast" works on TV shows
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_cast"] = strconv.Itoa(actorID)
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowShort(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetMoviesByDirector(directorID int, page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_crew"] = strconv.Itoa(directorID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) GetTVShowsByDirector(directorID int, page int) (*[]TVShow, error) {
	//TODO: Test this function, not sure "with_crew" works on TV shows
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_crew"] = strconv.Itoa(directorID)
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowShort(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetMoviesByStudio(studioID int, page int) (*[]Movie, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_companies"] = strconv.Itoa(studioID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = *extractMovieShort(&movie)
	}
	return &extractedMovies, nil
}

func (m *mediaClient) GetTVShowsByNetwork(studioID int, page int) (*[]TVShow, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_networks"] = strconv.Itoa(studioID)
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = *extractTVShowShort(&tvShow)
	}
	return &extractedTVShows, nil
}

func (m *mediaClient) GetTVShowsNextEpisode(tvIds []int, startDate, endDate string) (*[]TVEpisodeRelease, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mediaClient) GetMoviesReleases(movieIds []int, startDate, endDate string) (*[]MovieRelease, error) {
	//TODO implement me
	panic("implement me")
}

// extractMovie extracts movie information from a tmdb.Movie object and returns a Movie object.
// It uses the tmdb.MovieCredits object to extract actors, crew and studios.
func extractMovie(movie *tmdb.Movie, credits *tmdb.MovieCredits) *Movie {
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
	}
}

// extractMovieShort extracts movie information from a tmdb.MovieShort object and returns a Movie object.
func extractMovieShort(movie *tmdb.MovieShort) *Movie {
	return &Movie{
		ID:          movie.ID,
		BackdropURL: backdropImgURL(movie.BackdropPath),
		PosterURL:   posterImgURL(movie.PosterPath),
		Title:       movie.Title,
		Overview:    movie.Overview,
		ReleaseDate: movie.ReleaseDate,
		VoteAverage: movie.VoteAverage,
		VoteCount:   int(movie.VoteCount),
	}
}

// extractTVShow extracts TV show information from a tmdb.TVShow object and returns a TVShow object.
func extractTVEpisode(tvId int, episode *tmdb.TvEpisode) *TVEpisode {
	return &TVEpisode{
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

// extractTVShow extracts TV show information from a tmdb.TVShow object and returns a TVShow object.
func extractTVShow(tvShow *tmdb.TV, credits *tmdb.TvCredits) *TVShow {
	return &TVShow{
		ID:          tvShow.ID,
		Actors:      *extractTVActors(credits),
		BackdropURL: backdropImgURL(tvShow.BackdropPath),
		Crew:        *extractTVCrew(credits),
		Genres:      *extractGenres(&tvShow.Genres),
		Overview:    tvShow.Overview,
		PosterURL:   posterImgURL(tvShow.PosterPath),
		ReleaseDate: tvShow.FirstAirDate,
		Networks:    *extractStudios(&tvShow.Networks),
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
	}
}

// extractTVShowShort extracts TV show information from a tmdb.TVShowShort object and returns a TVShow object.
func extractTVShowShort(tvShow *tmdb.TvShort) *TVShow {
	return &TVShow{
		ID:          tvShow.ID,
		BackdropURL: backdropImgURL(tvShow.BackdropPath),
		PosterURL:   posterImgURL(tvShow.PosterPath),
		Title:       tvShow.Name,
		Overview:    tvShow.Overview,
		ReleaseDate: tvShow.FirstAirDate,
		VoteAverage: tvShow.VoteAverage,
		VoteCount:   int(tvShow.VoteCount),
	}
}

func extractTVShowResult(tvShow *struct {
	BackdropPath  string `json:"backdrop_path"`
	ID            int
	OriginalName  string   `json:"original_name"`
	FirstAirDate  string   `json:"first_air_date"`
	OriginCountry []string `json:"origin_country"`
	PosterPath    string   `json:"poster_path"`
	Popularity    float32
	Name          string
	VoteAverage   float32 `json:"vote_average"`
	VoteCount     uint32  `json:"vote_count"`
}) *TVShow {
	return &TVShow{
		ID:          tvShow.ID,
		BackdropURL: backdropImgURL(tvShow.BackdropPath),
		PosterURL:   posterImgURL(tvShow.PosterPath),
		Title:       tvShow.Name,
		ReleaseDate: tvShow.FirstAirDate,
		VoteAverage: tvShow.VoteAverage,
		VoteCount:   int(tvShow.VoteCount),
	}
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

func extractOptions(options map[string]string) map[string]string {
	var opts = make(map[string]string)
	for key, value := range options {
		opts[key] = value
	}
	return opts
}
