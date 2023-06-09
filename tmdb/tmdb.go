package tmdb

import (
	"fmt"
	"github.com/ryanbradynd05/go-tmdb"
	"log"
	"math"
	"sort"
	"strconv"
	"sync"
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
	ProfileURL string `json:"profileUrl"`
}

// Actor represents a person
type Actor struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProfileURL string `json:"profileUrl"`
	Overview   string `json:"overview"`
}

// Studio represents a movie/TV studio with its ID, name, and logo URL.
type Studio struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	LogoURL string `json:"logoUrl"`
}

// Movie represents a movie with its attributes such as ID, actors list (Person), backdrop URL,
// crew list (Person), genre list (Genre), overview, poster URL, release date, studio list (Studio),
// title, vote average, and vote count.
type Movie struct {
	ID          int      `json:"id"`
	Actors      []Person `json:"actors"`
	BackdropURL string   `json:"backdropUrl"`
	Crew        []Person `json:"crew"`
	Genres      []Genre  `json:"genres"`
	Overview    string   `json:"overview"`
	PosterURL   string   `json:"posterUrl"`
	ReleaseDate string   `json:"releaseDate"`
	Studios     []Studio `json:"studios"`
	Title       string   `json:"title"`
	VoteAverage float32  `json:"voteAverage"`
	VoteCount   int      `json:"voteCount"`
}

// TVEpisode represents a TV episode with its attributes such as ID, TV show ID, poster URL,
// season number, episode number, name, overview, and air date.
type TVEpisode struct {
	ID            int    `json:"id"`
	TVShowID      int    `json:"tvShowId"`
	PosterURL     string `json:"posterUrl"`
	EpisodeNumber int    `json:"episodeNumber"`
	SeasonNumber  int    `json:"seasonNumber"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	AirDate       string `json:"airDate"`
}

// TVShow represents a TV show with its attributes such as ID, actors list (Person), backdrop URL,
// crew list (Person), genre list (Genre), overview, poster URL, release date, studio list (Studio),
// status, next episode (TVEpisode), title, seasons count, vote average, and vote count.
type TVShow struct {
	ID            int        `json:"id"`
	Actors        []Person   `json:"actors"`
	BackdropURL   string     `json:"backdropUrl"`
	Crew          []Person   `json:"crew"`
	Genres        []Genre    `json:"genres"`
	Overview      string     `json:"overview"`
	PosterURL     string     `json:"posterUrl"`
	ReleaseDate   string     `json:"releaseDate"`
	Networks      []Studio   `json:"networks"`
	Status        string     `json:"status"`
	NextEpisode   *TVEpisode `json:"nextEpisode"`
	Title         string     `json:"title"`
	SeasonsCount  int        `json:"seasonsCount"`
	EpisodesCount int        `json:"episodesCount"`
	VoteAverage   float32    `json:"voteAverage"`
	VoteCount     int        `json:"voteCount"`
}

type PaginatedMovieResults struct {
	Results     []*Movie
	TotalPage   int
	TotalResult int
}

type PaginatedTVShowResults struct {
	Results     []*TVShow
	TotalPage   int
	TotalResult int
}

type PaginatedActorResults struct {
	Results     []*Actor
	TotalPage   int
	TotalResult int
}

// MediaClient is an interface for a media client API.
type MediaClient interface {
	GetActor(actorID int) (*Actor, error)
	GetMovie(id int) (*Movie, error)
	GetMovieGenre(genreID int) (*Genre, error)
	GetMovieGenres() ([]*Genre, error)
	GetMovieRecommendations(movieID int) ([]*Movie, error)
	GetMoviesByActor(actorID int, page int) (*PaginatedMovieResults, error)
	GetMoviesByDirector(directorID int, page int) (*PaginatedMovieResults, error)
	GetMoviesByGenre(genreID int, page int) (*PaginatedMovieResults, error)
	GetMoviesByStudio(studioID int, page int) (*PaginatedMovieResults, error)
	GetMovieShort(movieID int) (*Movie, error)
	GetMoviesReleases(movieIds []int, startDate, endDate time.Time) ([]*Movie, error)
	GetNetwork(networkID int) (*Studio, error)
	GetPopularMovies(page int) (*PaginatedMovieResults, error)
	GetPopularTVShows(page int) (*PaginatedTVShowResults, error)
	GetRecentMovies() ([]*Movie, error)
	GetRecentTVShows() ([]*TVShow, error)
	GetStudio(studioID int) (*Studio, error)
	GetTVEpisode(tvID, season, episodeNumber int) (*TVEpisode, error)
	GetTVGenre(genreID int) (*Genre, error)
	GetTVSeasonEpisodes(id int, season int) ([]*TVEpisode, error)
	GetTVShow(id int) (*TVShow, error)
	GetTVShowGenres() ([]*Genre, error)
	GetTVShowRecommendations(tvShowID int) ([]*TVShow, error)
	GetTVShowsByActor(actorID int, page int) (*PaginatedTVShowResults, error)
	GetTVShowsByGenre(genreID int, page int) (*PaginatedTVShowResults, error)
	GetTVShowsByNetwork(studioID int, page int) (*PaginatedTVShowResults, error)
	GetTVShowShort(tvShowID int) (*TVShow, error)
	GetTVShowsReleases(tvIds []int, startDate, endDate time.Time) ([]*TVEpisode, []*TVShow, error)
	SearchMovies(query string, page int, adult bool) (*PaginatedMovieResults, error)
	SearchMoviesYear(query string, year string, page int) (*PaginatedMovieResults, error)
	SearchTVShows(query string, page int, adult bool) (*PaginatedTVShowResults, error)
	SearchActors(query string, page int, adult bool) (*PaginatedActorResults, error)
}

type mediaClient struct {
	tmdbClient *tmdb.TMDb
	cache      mediaCache
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
			"region":   "fr",
		},
		cache: newInMemoryMediaCache(),
	}
}

func NewRedisMediaClient(apiKey, redisHost, redisPass string) MediaClient {
	config := tmdb.Config{
		APIKey:   apiKey,
		Proxies:  nil,
		UseProxy: false,
	}
	return &mediaClient{
		tmdbClient: tmdb.Init(config),
		options: map[string]string{
			"language": "fr",
			"region":   "fr",
		},
		cache: newRedisMediaCache(redisHost, redisPass),
	}
}

// GetMovie retrieves movie info and credits by ID and returns a Movie object.
func (m *mediaClient) GetMovie(id int) (*Movie, error) {
	cachedMovie := m.cache.GetMovie(id)
	if cachedMovie != nil {
		return cachedMovie, nil
	}

	movie, err := m.tmdbClient.GetMovieInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	m.cache.AddMovieShort(extractMovie(movie, nil))
	credits, err := m.tmdbClient.GetMovieCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	extracted := extractMovie(movie, credits)
	m.cache.AddMovie(extracted)

	return extracted, nil
}

// GetTVShow retrieves TV show info and credits by ID and returns a TVShow object.
func (m *mediaClient) GetTVShow(id int) (*TVShow, error) {
	cachedTVShow := m.cache.GetTV(id)
	if cachedTVShow != nil {
		return cachedTVShow, nil
	}

	tvShow, err := m.tmdbClient.GetTvInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	m.cache.AddTVShort(extractTVShow(tvShow, nil))
	credits, err := m.tmdbClient.GetTvCredits(id, m.options)
	if err != nil {
		return nil, err
	}
	extracted := extractTVShow(tvShow, credits)
	m.cache.AddTV(extracted)

	return extracted, nil
}

// GetMovieShort retrieves movie info by ID and returns a Movie object.
func (m *mediaClient) GetMovieShort(id int) (*Movie, error) {
	cachedMovie := m.cache.GetMovieShort(id)
	if cachedMovie != nil {
		return cachedMovie, nil
	}

	movie, err := m.tmdbClient.GetMovieInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	extracted := extractMovie(movie, nil)
	m.cache.AddMovieShort(extracted)

	return extracted, nil
}

// GetTVShowShort retrieves TV show info by ID and returns a TVShow object.
func (m *mediaClient) GetTVShowShort(id int) (*TVShow, error) {
	cachedTVShow := m.cache.GetTVShort(id)
	if cachedTVShow != nil {
		return cachedTVShow, nil
	}
	tvShow, err := m.tmdbClient.GetTvInfo(id, m.options)
	if err != nil {
		return nil, err
	}
	extracted := extractTVShow(tvShow, nil)
	m.cache.AddTVShort(extracted)

	return extracted, nil
}

// GetTVEpisode retrieves the information of a TV episode by TV show ID, season number and episode number and returns a TVEpisode object.
func (m *mediaClient) GetTVEpisode(tvID, season, episodeNumber int) (*TVEpisode, error) {
	cachedEpisode := m.cache.GetEpisode(tvID, season, episodeNumber)
	if cachedEpisode != nil {
		return cachedEpisode, nil
	}

	episode, err := m.tmdbClient.GetTvEpisodeInfo(tvID, season, episodeNumber, m.options)
	if err != nil {
		return nil, err
	}
	extracted := extractTVEpisode(tvID, episode)
	m.cache.AddEpisode(extracted)

	return extracted, nil
}

// GetTVSeasonEpisodes retrieves all episodes from a TV show season and returns a slice of TVEpisode objects.
func (m *mediaClient) GetTVSeasonEpisodes(tvID int, season int) ([]*TVEpisode, error) {
	cachedEpisodes := m.cache.GetSeason(tvID, season)
	if cachedEpisodes != nil {
		return cachedEpisodes, nil
	}

	episodes, err := m.tmdbClient.GetTvSeasonInfo(tvID, season, m.options)
	if err != nil {
		return nil, err
	}
	var extractedEpisodes = make([]*TVEpisode, len(episodes.Episodes))
	for i, episode := range episodes.Episodes {
		extractedEpisodes[i] = extractTVEpisode(tvID, &episode)
	}
	m.cache.AddSeason(tvID, season, extractedEpisodes)
	return extractedEpisodes, nil
}

// GetPopularMovies retrieves the most popular movies and returns a slice of Movie objects.
func (m *mediaClient) GetPopularMovies(page int) (*PaginatedMovieResults, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	movies, err := m.tmdbClient.GetMoviePopular(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	return &PaginatedMovieResults{
		Results:     extractedMovies,
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
	}, nil
}

// GetPopularTVShows retrieves the most popular TV shows and returns a slice of TVShow objects.
func (m *mediaClient) GetPopularTVShows(page int) (*PaginatedTVShowResults, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	tvShows, err := m.tmdbClient.GetTvPopular(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]*TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = extractTVShowShort(&tvShow)
	}
	return &PaginatedTVShowResults{
		Results:     extractedTVShows,
		TotalPage:   tvShows.TotalPages,
		TotalResult: tvShows.TotalResults,
	}, nil
}

// GetRecentMovies retrieves the most recent movies and returns a slice of Movie objects.
func (m *mediaClient) GetRecentMovies() ([]*Movie, error) {
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
	var extractedMovies = make([]*Movie, 0)
	// Get the 20 most popular
	for _, movie := range movies[:20] {
		extractedMovies = append(extractedMovies, extractMovieShort(&movie))
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
	return extractedMovies, nil
}

// GetRecentTVShows retrieves the most recent TV shows and returns a slice of TVShow objects.
func (m *mediaClient) GetRecentTVShows() ([]*TVShow, error) {
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
	var extractedTVShows = make([]*TVShow, 0)
	// Get the 20 most popular
	for _, tvshow := range tvshows[:20] {
		extractedTVShows = append(extractedTVShows, extractTVShowShort(&tvshow))
	}
	// Return result
	return extractedTVShows, nil
}

// SearchMovies searches for movies matching the given query and returns a slice of Movie objects.
func (m *mediaClient) SearchMovies(query string, page int, adult bool) (*PaginatedMovieResults, error) {
	cachedResults := m.cache.GetMovieSearchResults(query, page, adult)
	if cachedResults != nil {
		return cachedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["region"] = "fr"
	if adult {
		options["include_adult"] = "true"
	}
	movies, err := m.tmdbClient.SearchMovie(query, options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	result := &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}
	m.cache.AddMovieSearchResults(query, page, adult, result)
	return result, nil
}

// SearchMoviesYear searches for movies matching the given query and year and returns a slice of Movie objects.
func (m *mediaClient) SearchMoviesYear(query string, year string, page int) (*PaginatedMovieResults, error) {
	cachedResults := m.cache.GetMovieSearchResultsYear(query, page, year)
	if cachedResults != nil {
		return cachedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["region"] = "fr"
	options["year"] = year
	movies, err := m.tmdbClient.SearchMovie(query, options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	result := &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}
	m.cache.AddMovieSearchResultsYear(query, page, year, result)
	return result, nil
}

// SearchTVShows searches for TV shows matching the given query and returns a slice of TVShow objects.
func (m *mediaClient) SearchTVShows(query string, page int, adult bool) (*PaginatedTVShowResults, error) {
	extractedResults := m.cache.GetTVSearchResults(query, page, adult)
	if extractedResults != nil {
		return extractedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	if adult {
		options["include_adult"] = "true"
	}
	tvShows, err := m.tmdbClient.SearchTv(query, options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]*TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = extractTVShowResult(&tvShow)
	}
	result := &PaginatedTVShowResults{
		TotalPage:   tvShows.TotalPages,
		TotalResult: tvShows.TotalResults,
		Results:     extractedTVShows,
	}
	m.cache.AddTVSearchResults(query, page, adult, result)
	return result, nil
}

// SearchActors searches for actors matching the given query and returns a slice of Actor objects.
func (m *mediaClient) SearchActors(query string, page int, adult bool) (*PaginatedActorResults, error) {
	extractedResults := m.cache.GetActorSearchResults(query, page, adult)
	if extractedResults != nil {
		return extractedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	if adult {
		options["include_adult"] = "true"
	}
	actors, err := m.tmdbClient.SearchPerson(query, options)
	if err != nil {
		return nil, err
	}
	var extractedActors = extractActors(actors.Results)

	result := &PaginatedActorResults{
		TotalPage:   actors.TotalPages,
		TotalResult: actors.TotalResults,
		Results:     extractedActors,
	}
	m.cache.AddActorSearchResults(query, page, adult, result)
	return result, nil
}

// GetMoviesByGenre retrieves movies of the given genre and returns a slice of Movie objects.
func (m *mediaClient) GetMoviesByGenre(genreID int, page int) (*PaginatedMovieResults, error) {
	cachedResults := m.cache.GetMoviesByGenre(genreID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_genres"] = strconv.Itoa(genreID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	result := &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}
	m.cache.AddMoviesByGenre(genreID, page, result)
	return result, nil
}

// GetTVShowsByGenre retrieves TV shows of the given genre and returns a slice of TVShow objects.
func (m *mediaClient) GetTVShowsByGenre(genreID int, page int) (*PaginatedTVShowResults, error) {
	cachedResults := m.cache.GetTVsByGenre(genreID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}

	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_genres"] = strconv.Itoa(genreID)
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]*TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = extractTVShowShort(&tvShow)
	}
	result := &PaginatedTVShowResults{
		TotalPage:   tvShows.TotalPages,
		TotalResult: tvShows.TotalResults,
		Results:     extractedTVShows,
	}
	m.cache.AddTVsByGenre(genreID, page, result)
	return result, nil
}

// GetMoviesByActor retrieves movies starring the given actor and returns a slice of Movie objects.
func (m *mediaClient) GetMoviesByActor(actorID int, page int) (*PaginatedMovieResults, error) {
	cachedResults := m.cache.GetMoviesByActor(actorID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_cast"] = strconv.Itoa(actorID)
	options["include_adult"] = "true"
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	result := &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}
	m.cache.AddMoviesByActor(actorID, page, result)
	return result, nil
}

func (m *mediaClient) GetTVShowsByActor(actorID int, page int) (*PaginatedTVShowResults, error) {
	cachedResults := m.cache.GetTVsByActor(actorID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}

	actorTVCredits, err := m.tmdbClient.GetPersonTvCredits(actorID, m.options)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var startIndex = int(math.Min(float64((page-1)*20), math.Max(0, float64(len(actorTVCredits.Cast)-1))))
	var endIndex = int(math.Min(float64(page*20), float64(len(actorTVCredits.Cast))))
	var extractedTVShows = make([]*TVShow, endIndex-startIndex)
	var lockIndexes = make([]sync.Mutex, endIndex-startIndex)
	for index, tvShow := range actorTVCredits.Cast[startIndex:endIndex] {
		wg.Add(1)
		go func(tvShowID, index int) {
			defer wg.Done()
			tvShow, err := m.GetTVShowShort(tvShowID)
			if err != nil {
				log.Printf("Error while retrieving TV show %d: %s", tvShowID, err)
				return
			}
			lockIndexes[index].Lock()
			defer lockIndexes[index].Unlock()
			extractedTVShows[index] = tvShow
		}(tvShow.ID, index)
	}
	wg.Wait()

	result := &PaginatedTVShowResults{
		TotalPage:   int(math.Round(float64(len(actorTVCredits.Cast)) / 20)),
		TotalResult: len(actorTVCredits.Cast),
		Results:     extractedTVShows,
	}
	m.cache.AddTVsByActor(actorID, page, result)
	return result, nil
}

// GetMoviesByDirector retrieves movies directed by the given director and returns a slice of Movie objects.
func (m *mediaClient) GetMoviesByDirector(directorID int, page int) (*PaginatedMovieResults, error) {
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_crew"] = strconv.Itoa(directorID)
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	return &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}, nil
}

// GetMoviesByStudio retrieves movies produced by the given studio and returns a slice of Movie objects.
func (m *mediaClient) GetMoviesByStudio(studioID int, page int) (*PaginatedMovieResults, error) {
	cachedResults := m.cache.GetMoviesByStudio(studioID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_companies"] = strconv.Itoa(studioID)
	options["include_adult"] = "true"
	movies, err := m.tmdbClient.DiscoverMovie(options)
	if err != nil {
		return nil, err
	}
	var extractedMovies = make([]*Movie, len(movies.Results))
	for i, movie := range movies.Results {
		extractedMovies[i] = extractMovieShort(&movie)
	}
	result := &PaginatedMovieResults{
		TotalPage:   movies.TotalPages,
		TotalResult: movies.TotalResults,
		Results:     extractedMovies,
	}
	m.cache.AddMoviesByStudio(studioID, page, result)
	return result, nil
}

// GetTVShowsByNetwork retrieves TV shows produced by the given studio and returns a slice of TVShow objects.
func (m *mediaClient) GetTVShowsByNetwork(studioID int, page int) (*PaginatedTVShowResults, error) {
	cachedResults := m.cache.GetTVsByNetwork(studioID, page)
	if cachedResults != nil {
		return cachedResults, nil
	}
	options := extractOptions(m.options)
	options["page"] = strconv.Itoa(page)
	options["with_networks"] = strconv.Itoa(studioID)
	options["include_adult"] = "true"
	tvShows, err := m.tmdbClient.DiscoverTV(options)
	if err != nil {
		return nil, err
	}
	var extractedTVShows = make([]*TVShow, len(tvShows.Results))
	for i, tvShow := range tvShows.Results {
		extractedTVShows[i] = extractTVShowShort(&tvShow)
	}
	result := &PaginatedTVShowResults{
		TotalPage:   tvShows.TotalPages,
		TotalResult: tvShows.TotalResults,
		Results:     extractedTVShows,
	}
	m.cache.AddTVsByNetwork(studioID, page, result)
	return result, nil
}

// GetTVShowsReleases retrieves all TV shows airing between the given dates and returns a slice of TVEpisodeRelease objects.
func (m *mediaClient) GetTVShowsReleases(tvIds []int, startDate, endDate time.Time) ([]*TVEpisode, []*TVShow, error) {
	// Get all episodes for the given TV shows that are airing between the given dates
	var episodes []*TVEpisode
	var tvShows []*TVShow
	var lock sync.Mutex
	var wg sync.WaitGroup
	for _, tvID := range tvIds {
		wg.Add(1)
		go func(tvID int) {
			defer wg.Done()
			tvShow, err := m.GetTVShowShort(tvID)
			if err != nil {
				log.Printf("Error while retrieving TV show %d: %s", tvID, err)
				return
			}
			// Get all episodes for the given TV show that are airing between the given dates
			showAdded := false
			for seasonNumber := 1; seasonNumber <= tvShow.SeasonsCount; seasonNumber++ {
				wg.Add(1)
				go func(tvID, seasonNumber int) {
					defer wg.Done()
					seasonEpisodes, err := m.GetTVSeasonEpisodes(tvID, seasonNumber)
					if err != nil {
						log.Printf("Error while retrieving TV show %d season %d: %s", tvID, seasonNumber, err)
						return
					}
					var episodesToAdd []*TVEpisode
					for _, episode := range seasonEpisodes {
						airDate, err := time.Parse("2006-01-02", episode.AirDate)
						if err != nil {
							log.Printf("Could not parse air date %s for episode %d of TV show %d",
								episode.AirDate, episode.ID, tvID)
							continue
						}
						if (airDate.After(startDate) && airDate.Before(endDate)) ||
							airDate.Equal(startDate) ||
							airDate.Equal(endDate) {
							episodesToAdd = append(episodesToAdd, episode)
						}
					}
					if len(episodesToAdd) > 0 {
						lock.Lock()
						defer lock.Unlock()
						episodes = append(episodes, episodesToAdd...)
						if !showAdded {
							tvShows = append(tvShows, tvShow)
							showAdded = true
						}
					}
				}(tvID, seasonNumber)
			}
		}(tvID)
	}
	wg.Wait()
	return episodes, tvShows, nil
}

// GetMoviesReleases retrieves all movies released between the given dates and returns a slice of MovieRelease objects.
func (m *mediaClient) GetMoviesReleases(movieIds []int, startDate, endDate time.Time) ([]*Movie, error) {
	var movies []*Movie
	var lock sync.Mutex
	var wg sync.WaitGroup
	for _, movieID := range movieIds {
		wg.Add(1)
		go func(movieID int) {
			defer wg.Done()
			movie, err := m.GetMovieShort(movieID)
			if err != nil {
				log.Printf("Error while retrieving movie %d: %s", movieID, err)
				return
			}
			airDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
			if err != nil {
				log.Printf("Could not parse air date %s for movie %d",
					movie.ReleaseDate, movie.ID)
				return
			}
			if (airDate.After(startDate) && airDate.Before(endDate)) ||
				airDate.Equal(startDate) ||
				airDate.Equal(endDate) {
				lock.Lock()
				defer lock.Unlock()
				movies = append(movies, movie)
			}
		}(movieID)
	}
	wg.Wait()
	return movies, nil
}

// GetMovieRecommendations retrieves movie recommendations for the given movie and returns a slice of Movie objects.
func (m *mediaClient) GetMovieRecommendations(movieID int) ([]*Movie, error) {
	cachedResults := m.cache.GetMovieRecommendations(movieID)
	if cachedResults != nil {
		return cachedResults, nil
	}
	recommendations, err := m.tmdbClient.GetMovieRecommendations(movieID, m.options)
	if err != nil {
		return nil, err
	}
	movies := make([]*Movie, len(recommendations.Results))
	for i, movieRecommendation := range recommendations.Results {
		movies[i] = extractMovieShort(&tmdb.MovieShort{
			ID:           movieRecommendation.ID,
			Title:        movieRecommendation.Title,
			Overview:     movieRecommendation.Overview,
			ReleaseDate:  movieRecommendation.ReleaseDate,
			PosterPath:   movieRecommendation.PosterPath,
			BackdropPath: movieRecommendation.BackdropPath,
			VoteAverage:  movieRecommendation.VoteAverage,
			VoteCount:    movieRecommendation.VoteCount,
		})
	}
	m.cache.AddMovieRecommendations(movieID, movies)
	return movies, nil
}

// GetTVShowRecommendations retrieves TV show recommendations for the given TV show and returns a slice of TVShow objects.
func (m *mediaClient) GetTVShowRecommendations(tvShowID int) ([]*TVShow, error) {
	cachedResults := m.cache.GetTVRecommendations(tvShowID)
	if cachedResults != nil {
		return cachedResults, nil
	}
	recommendations, err := m.tmdbClient.GetTvRecommendations(tvShowID, m.options)
	if err != nil {
		return nil, err
	}
	tvShows := make([]*TVShow, len(recommendations.Results))
	for i, tvShowRecommendation := range recommendations.Results {
		tvShows[i] = extractTVShowShort(&tmdb.TvShort{
			ID:           tvShowRecommendation.ID,
			Name:         tvShowRecommendation.Name,
			Overview:     tvShowRecommendation.Overview,
			FirstAirDate: tvShowRecommendation.FirstAirDate,
			PosterPath:   tvShowRecommendation.PosterPath,
			BackdropPath: tvShowRecommendation.BackdropPath,
			VoteAverage:  tvShowRecommendation.VoteAverage,
			VoteCount:    tvShowRecommendation.VoteCount,
		})
	}
	m.cache.AddTVRecommendations(tvShowID, tvShows)
	return tvShows, nil
}

func (m *mediaClient) GetMovieGenre(genreID int) (*Genre, error) {
	cachedGenre := m.cache.GetMovieGenre(genreID)
	if cachedGenre != nil {
		return cachedGenre, nil
	}

	genres, err := m.tmdbClient.GetMovieGenres(m.options)
	if err != nil {
		return nil, err
	}
	for _, genre := range genres.Genres {
		if genre.ID == genreID {
			genre := &Genre{
				ID:   genre.ID,
				Name: genre.Name,
			}
			m.cache.AddMovieGenre(genre)
			return genre, nil
		}
	}
	return nil, fmt.Errorf("movie genre with ID %d not found", genreID)
}

func (m *mediaClient) GetTVGenre(genreID int) (*Genre, error) {
	cachedGenre := m.cache.GetTVGenre(genreID)
	if cachedGenre != nil {
		return cachedGenre, nil
	}

	genres, err := m.tmdbClient.GetTvGenres(m.options)
	if err != nil {
		return nil, err
	}
	for _, genre := range genres.Genres {
		if genre.ID == genreID {
			genre := &Genre{
				ID:   genre.ID,
				Name: genre.Name,
			}
			m.cache.AddTVGenre(genre)
			return genre, nil
		}
	}
	return nil, fmt.Errorf("TV genre with ID %d not found", genreID)
}

func (m *mediaClient) GetMovieGenres() ([]*Genre, error) {
	genres, err := m.tmdbClient.GetMovieGenres(m.options)
	if err != nil {
		return nil, err
	}
	movieGenres := make([]*Genre, len(genres.Genres))
	for i, genre := range genres.Genres {
		movieGenres[i] = &Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}
	}
	return movieGenres, nil
}

func (m *mediaClient) GetTVShowGenres() ([]*Genre, error) {
	genres, err := m.tmdbClient.GetTvGenres(m.options)
	if err != nil {
		return nil, err
	}
	tvGenres := make([]*Genre, len(genres.Genres))
	for i, genre := range genres.Genres {
		tvGenres[i] = &Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}
	}
	return tvGenres, nil
}

func (m *mediaClient) GetActor(actorID int) (*Actor, error) {
	cachedActor := m.cache.GetActor(actorID)
	if cachedActor != nil {
		return cachedActor, nil
	}

	response, err := m.tmdbClient.GetPersonInfo(actorID, m.options)
	if err != nil {
		return nil, err
	}
	actor := &Actor{
		ID:         response.ID,
		Name:       response.Name,
		ProfileURL: profileImgURL(response.ProfilePath),
		Overview:   response.Biography,
	}
	m.cache.AddActor(actor)
	return actor, nil
}

func (m *mediaClient) GetStudio(studioID int) (*Studio, error) {
	response, err := m.tmdbClient.GetCompanyInfo(studioID, m.options)
	if err != nil {
		return nil, err
	}
	return &Studio{
		ID:      response.ID,
		Name:    response.Name,
		LogoURL: profileImgURL(response.LogoPath),
	}, nil
}

func (m *mediaClient) GetNetwork(networkID int) (*Studio, error) {
	response, err := m.tmdbClient.GetNetworkInfo(networkID)
	if err != nil {
		return nil, err
	}
	return &Studio{
		ID:      response.ID,
		Name:    response.Name,
		LogoURL: profileImgURL(""),
	}, nil
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
		PosterURL:     backdropImgURL(episode.StillPath),
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
				PosterURL:     backdropImgURL(tvShow.NextEpisodeToAir.StillPath),
				EpisodeNumber: tvShow.NextEpisodeToAir.EpisodeNumber,
				SeasonNumber:  tvShow.NextEpisodeToAir.SeasonNumber,
				Name:          tvShow.NextEpisodeToAir.Name,
				Overview:      tvShow.NextEpisodeToAir.Overview,
				AirDate:       tvShow.NextEpisodeToAir.AirDate,
			}
		}(),
		SeasonsCount:  tvShow.NumberOfSeasons,
		EpisodesCount: tvShow.NumberOfEpisodes,
		VoteAverage:   tvShow.VoteAverage,
		VoteCount:     int(tvShow.VoteCount),
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
	if credits == nil {
		return &[]Person{}
	}
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
	if credits == nil {
		return &[]Person{}
	}
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

// extractActors extracts actors from credits and returns a list of Person.
func extractActors(actors []struct {
	Adult       bool
	ID          int
	Name        string
	Popularity  float32
	ProfilePath string `json:"profile_path"`
	KnownFor    []struct {
		Adult         bool
		BackdropPath  string `json:"backdrop_path"`
		ID            int
		OriginalTitle string `json:"original_title"`
		ReleaseDate   string `json:"release_date"`
		PosterPath    string `json:"poster_path"`
		Popularity    float32
		Title         string
		VoteAverage   float32 `json:"vote_average"`
		VoteCount     uint32  `json:"vote_count"`
		MediaType     string  `json:"media_type"`
	} `json:"known_for"`
}) []*Actor {
	var cast = make([]*Actor, len(actors))
	for i, actor := range actors {
		cast[i] = &Actor{
			ID:         actor.ID,
			Name:       actor.Name,
			ProfileURL: profileImgURL(actor.ProfilePath),
		}
	}
	return cast
}

// extractMovieCrew extracts crew from movie credits and returns a list of Person.
func extractMovieCrew(credits *tmdb.MovieCredits) *[]Person {
	if credits == nil {
		return &[]Person{}
	}
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
	if credits == nil {
		return &[]Person{}
	}
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
