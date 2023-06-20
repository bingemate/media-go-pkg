package tmdb

import (
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
	"log"
	"strconv"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type mediaCache interface {
	AddMovie(m *Movie)
	GetMovie(id int) *Movie
	AddMovieShort(m *Movie)
	GetMovieShort(id int) *Movie
	AddTV(t *TVShow)
	GetTV(id int) *TVShow
	AddTVShort(t *TVShow)
	GetTVShort(id int) *TVShow
	AddEpisode(e *TVEpisode)
	GetEpisode(tvID int, seasonNumber int, episodeNumber int) *TVEpisode
	AddSeason(tvID int, seasonNumber int, s []*TVEpisode)
	GetSeason(tvID int, seasonNumber int) []*TVEpisode
	AddMovieSearchResults(query string, page int, results *PaginatedMovieResults)
	GetMovieSearchResults(query string, page int) *PaginatedMovieResults
	GetMovieSearchResultsYear(query string, page int, year string) *PaginatedMovieResults
	AddMovieSearchResultsYear(query string, page int, year string, results *PaginatedMovieResults)
	AddTVSearchResults(query string, page int, results *PaginatedTVShowResults)
	GetTVSearchResults(query string, page int) *PaginatedTVShowResults
	AddMovieGenre(genre *Genre)
	GetMovieGenre(id int) *Genre
	AddTVGenre(genre *Genre)
	GetTVGenre(id int) *Genre
	AddActor(actor *Actor)
	GetActor(id int) *Actor
}

type inMemoryMediaCache struct {
	cache *cache.Cache
}

func newInMemoryMediaCache() mediaCache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &inMemoryMediaCache{
		cache: c,
	}
}

func (c *inMemoryMediaCache) AddMovie(m *Movie) {
	c.cache.SetDefault("movie:"+strconv.Itoa(m.ID), m)
}

func (c *inMemoryMediaCache) GetMovie(id int) *Movie {
	m, ok := c.cache.Get("movie:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return m.(*Movie)
}

func (c *inMemoryMediaCache) AddMovieShort(m *Movie) {
	c.cache.SetDefault("movie_short:"+strconv.Itoa(m.ID), m)
}

func (c *inMemoryMediaCache) GetMovieShort(id int) *Movie {
	m, ok := c.cache.Get("movie_short:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return m.(*Movie)
}

func (c *inMemoryMediaCache) AddTV(t *TVShow) {
	c.cache.SetDefault("tv:"+strconv.Itoa(t.ID), t)
}

func (c *inMemoryMediaCache) GetTV(id int) *TVShow {
	t, ok := c.cache.Get("tv:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return t.(*TVShow)
}

func (c *inMemoryMediaCache) AddTVShort(t *TVShow) {
	c.cache.SetDefault("tv_short:"+strconv.Itoa(t.ID), t)
}

func (c *inMemoryMediaCache) GetTVShort(id int) *TVShow {
	t, ok := c.cache.Get("tv_short:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return t.(*TVShow)
}

func (c *inMemoryMediaCache) AddEpisode(e *TVEpisode) {
	c.cache.SetDefault("episode:"+strconv.Itoa(e.TVShowID)+":"+strconv.Itoa(e.SeasonNumber)+":"+strconv.Itoa(e.EpisodeNumber), e)
}

func (c *inMemoryMediaCache) GetEpisode(tvID int, seasonNumber int, episodeNumber int) *TVEpisode {
	e, ok := c.cache.Get("episode:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber) + ":" + strconv.Itoa(episodeNumber))
	if !ok {
		return nil
	}
	return e.(*TVEpisode)
}

func (c *inMemoryMediaCache) AddSeason(tvID int, seasonNumber int, s []*TVEpisode) {
	c.cache.SetDefault("season:"+strconv.Itoa(tvID)+":"+strconv.Itoa(seasonNumber), s)
}

func (c *inMemoryMediaCache) GetSeason(tvID int, seasonNumber int) []*TVEpisode {
	s, ok := c.cache.Get("season:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber))
	if !ok {
		return nil
	}
	return s.([]*TVEpisode)
}

func (c *inMemoryMediaCache) AddMovieSearchResults(query string, page int, results *PaginatedMovieResults) {
	c.cache.SetDefault("movie_search:"+query+":"+strconv.Itoa(page), results)
}

func (c *inMemoryMediaCache) GetMovieSearchResults(query string, page int) *PaginatedMovieResults {
	r, ok := c.cache.Get("movie_search:" + query + ":" + strconv.Itoa(page))
	if !ok {
		return nil
	}
	return r.(*PaginatedMovieResults)
}

func (c *inMemoryMediaCache) GetMovieSearchResultsYear(query string, page int, year string) *PaginatedMovieResults {
	r, ok := c.cache.Get("movie_search:" + query + ":" + strconv.Itoa(page) + ":" + year)
	if !ok {
		return nil
	}
	return r.(*PaginatedMovieResults)
}

func (c *inMemoryMediaCache) AddMovieSearchResultsYear(query string, page int, year string, results *PaginatedMovieResults) {
	c.cache.SetDefault("movie_search:"+query+":"+strconv.Itoa(page)+":"+year, results)
}

func (c *inMemoryMediaCache) AddTVSearchResults(query string, page int, results *PaginatedTVShowResults) {
	c.cache.SetDefault("tv_search:"+query+":"+strconv.Itoa(page), results)
}

func (c *inMemoryMediaCache) GetTVSearchResults(query string, page int) *PaginatedTVShowResults {
	r, ok := c.cache.Get("tv_search:" + query + ":" + strconv.Itoa(page))
	if !ok {
		return nil
	}
	return r.(*PaginatedTVShowResults)
}

func (c *inMemoryMediaCache) AddMovieGenre(genre *Genre) {
	c.cache.SetDefault("movie_genre:"+strconv.Itoa(genre.ID), genre)
}

func (c *inMemoryMediaCache) GetMovieGenre(id int) *Genre {
	g, ok := c.cache.Get("movie_genre:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return g.(*Genre)
}

func (c *inMemoryMediaCache) AddTVGenre(genre *Genre) {
	c.cache.SetDefault("tv_genre:"+strconv.Itoa(genre.ID), genre)
}

func (c *inMemoryMediaCache) GetTVGenre(id int) *Genre {
	g, ok := c.cache.Get("tv_genre:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return g.(*Genre)
}

func (c *inMemoryMediaCache) AddActor(actor *Actor) {
	c.cache.SetDefault("actor:"+strconv.Itoa(actor.ID), actor)
}

func (c *inMemoryMediaCache) GetActor(id int) *Actor {
	a, ok := c.cache.Get("actor:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return a.(*Actor)
}

type redisMediaCache struct {
	client *redis.Client
}

func newRedisMediaCache(redisURL string, redisPassword string) mediaCache {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       0,
	})
	return &redisMediaCache{
		client: client,
	}
}

var (
	defaultExpiration = 30 * 24 * time.Hour // 1 mois
	oneWeekExpiration = 7 * 24 * time.Hour  // 1 semaine
)

/*
- Si le film / la série / l'épisode a une date de moins de 1 mois -> 1 semaine de rétention
  Sinon, rétention de 1 mois
- Saison -> Rétention 1 semaine
- Résultat de recherche film / série -> 1 semaine de rétention
- Genre et Acteur -> 1 mois rétention
*/

func calculateExpirationDate(releaseDate string, defaultExpiration, recentExpiration time.Duration) time.Duration {
	if releaseDate == "" {
		return defaultExpiration
	}
	releaseDateParsed, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return defaultExpiration
	}

	diff := time.Now().Sub(releaseDateParsed)
	if diff < 30*24*time.Hour {
		return recentExpiration
	}
	return defaultExpiration
}

func (r redisMediaCache) AddMovie(m *Movie) {
	key := "movie:" + strconv.Itoa(m.ID)
	expiration := calculateExpirationDate(m.ReleaseDate, defaultExpiration, oneWeekExpiration)

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error while marshalling movie", err)
		return
	}
	r.client.Set(key, data, expiration)
}

func (r redisMediaCache) GetMovie(id int) *Movie {
	key := "movie:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var m Movie
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Println("Error while unmarshalling movie", err)
		return nil
	}
	return &m
}

func (r redisMediaCache) AddMovieShort(m *Movie) {
	key := "movie_short:" + strconv.Itoa(m.ID)
	expiration := calculateExpirationDate(m.ReleaseDate, defaultExpiration, oneWeekExpiration)

	data, err := json.Marshal(m)
	if err != nil {
		log.Println("Error while marshalling movie short", err)
		return
	}
	r.client.Set(key, data, expiration)
}

func (r redisMediaCache) GetMovieShort(id int) *Movie {
	key := "movie_short:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var m Movie
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Println("Error while unmarshalling movie short", err)
		return nil
	}
	return &m
}

func (r redisMediaCache) AddTV(t *TVShow) {
	key := "tv:" + strconv.Itoa(t.ID)
	expiration := calculateExpirationDate(t.ReleaseDate, defaultExpiration, oneWeekExpiration)

	data, err := json.Marshal(t)
	if err != nil {
		log.Println("Error while marshalling tv show", err)
		return
	}
	r.client.Set(key, data, expiration)
}

func (r redisMediaCache) GetTV(id int) *TVShow {
	key := "tv:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var t TVShow
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Println("Error while unmarshalling tv show", err)
		return nil
	}
	return &t
}

func (r redisMediaCache) AddTVShort(t *TVShow) {
	key := "tv_short:" + strconv.Itoa(t.ID)
	expiration := calculateExpirationDate(t.ReleaseDate, defaultExpiration, oneWeekExpiration)

	data, err := json.Marshal(t)
	if err != nil {
		log.Println("Error while marshalling tv show short", err)
		return
	}
	r.client.Set(key, data, expiration)
}

func (r redisMediaCache) GetTVShort(id int) *TVShow {
	key := "tv_short:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var t TVShow
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Println("Error while unmarshalling tv show short", err)
		return nil
	}
	return &t
}

func (r redisMediaCache) AddEpisode(e *TVEpisode) {
	key := "episode:" + strconv.Itoa(e.TVShowID) + ":" + strconv.Itoa(e.SeasonNumber) + ":" + strconv.Itoa(e.EpisodeNumber)
	expiration := calculateExpirationDate(e.AirDate, defaultExpiration, oneWeekExpiration)

	data, err := json.Marshal(e)
	if err != nil {
		log.Println("Error while marshalling episode", err)
		return
	}
	r.client.Set(key, data, expiration)
}

func (r redisMediaCache) GetEpisode(tvID int, seasonNumber int, episodeNumber int) *TVEpisode {
	key := "episode:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber) + ":" + strconv.Itoa(episodeNumber)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var e TVEpisode
	err = json.Unmarshal(data, &e)
	if err != nil {
		log.Println("Error while unmarshalling episode", err)
		return nil
	}
	return &e
}

func (r redisMediaCache) AddSeason(tvID int, seasonNumber int, s []*TVEpisode) {
	key := "season:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber)
	data, err := json.Marshal(s)
	if err != nil {
		log.Println("Error while marshalling season", err)
		return
	}
	r.client.Set(key, data, defaultExpiration)
	for _, e := range s {
		r.AddEpisode(e)
	}
}

func (r redisMediaCache) GetSeason(tvID int, seasonNumber int) []*TVEpisode {
	key := "season:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var s []*TVEpisode
	err = json.Unmarshal(data, &s)
	if err != nil {
		log.Println("Error while unmarshalling season", err)
		return nil
	}
	return s
}

func (r redisMediaCache) AddMovieSearchResults(query string, page int, results *PaginatedMovieResults) {
	key := "movie_search:" + query + ":" + strconv.Itoa(page)
	data, err := json.Marshal(results)
	if err != nil {
		log.Println("Error while marshalling movie search results", err)
		return
	}
	r.client.Set(key, data, oneWeekExpiration)
}

func (r redisMediaCache) GetMovieSearchResults(query string, page int) *PaginatedMovieResults {
	key := "movie_search:" + query + ":" + strconv.Itoa(page)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var results PaginatedMovieResults
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Println("Error while unmarshalling movie search results", err)
		return nil
	}
	return &results
}

func (r redisMediaCache) GetMovieSearchResultsYear(query string, page int, year string) *PaginatedMovieResults {
	key := "movie_search:" + query + ":" + strconv.Itoa(page) + ":" + year
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var results PaginatedMovieResults
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Println("Error while unmarshalling movie search results", err)
		return nil
	}
	return &results
}

func (r redisMediaCache) AddMovieSearchResultsYear(query string, page int, year string, results *PaginatedMovieResults) {
	key := "movie_search:" + query + ":" + strconv.Itoa(page) + ":" + year
	data, err := json.Marshal(results)
	if err != nil {
		log.Println("Error while marshalling movie search results", err)
		return
	}
	r.client.Set(key, data, oneWeekExpiration)
}

func (r redisMediaCache) AddTVSearchResults(query string, page int, results *PaginatedTVShowResults) {
	key := "tv_search:" + query + ":" + strconv.Itoa(page)
	data, err := json.Marshal(results)
	if err != nil {
		log.Println("Error while marshalling tv search results", err)
		return
	}
	r.client.Set(key, data, oneWeekExpiration)
}

func (r redisMediaCache) GetTVSearchResults(query string, page int) *PaginatedTVShowResults {
	key := "tv_search:" + query + ":" + strconv.Itoa(page)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var results PaginatedTVShowResults
	err = json.Unmarshal(data, &results)
	if err != nil {
		log.Println("Error while unmarshalling tv search results", err)
		return nil
	}
	return &results
}

func (r redisMediaCache) AddMovieGenre(genre *Genre) {
	key := "movie_genre:" + strconv.Itoa(genre.ID)
	data, err := json.Marshal(genre)
	if err != nil {
		log.Println("Error while marshalling movie genre", err)
		return
	}
	r.client.Set(key, data, defaultExpiration)
}

func (r redisMediaCache) GetMovieGenre(id int) *Genre {
	key := "movie_genre:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var g Genre
	err = json.Unmarshal(data, &g)
	if err != nil {
		log.Println("Error while unmarshalling movie genre", err)
		return nil
	}
	return &g
}

func (r redisMediaCache) AddTVGenre(genre *Genre) {
	key := "tv_genre:" + strconv.Itoa(genre.ID)
	data, err := json.Marshal(genre)
	if err != nil {
		log.Println("Error while marshalling tv genre", err)
		return
	}
	r.client.Set(key, data, defaultExpiration)
}

func (r redisMediaCache) GetTVGenre(id int) *Genre {
	key := "tv_genre:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var g Genre
	err = json.Unmarshal(data, &g)
	if err != nil {
		log.Println("Error while unmarshalling tv genre", err)
		return nil
	}
	return &g
}

func (r redisMediaCache) AddActor(actor *Actor) {
	key := "actor:" + strconv.Itoa(actor.ID)
	data, err := json.Marshal(actor)
	if err != nil {
		log.Println("Error while marshalling actor", err)
		return
	}
	r.client.Set(key, data, defaultExpiration)
}

func (r redisMediaCache) GetActor(id int) *Actor {
	key := "actor:" + strconv.Itoa(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		return nil
	}
	var a Actor
	err = json.Unmarshal(data, &a)
	if err != nil {
		log.Println("Error while unmarshalling actor", err)
		return nil
	}
	return &a
}
