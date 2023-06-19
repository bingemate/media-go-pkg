package tmdb

import (
	"github.com/patrickmn/go-cache"
	"strconv"
	"time"
)

type mediaCache struct {
	cache *cache.Cache
}

func newMediaCache() *mediaCache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &mediaCache{
		cache: c,
	}
}

func (c *mediaCache) AddMovie(m *Movie) {
	c.cache.SetDefault("movie:"+strconv.Itoa(m.ID), m)
}

func (c *mediaCache) GetMovie(id int) *Movie {
	m, ok := c.cache.Get("movie:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return m.(*Movie)
}

func (c *mediaCache) AddMovieShort(m *Movie) {
	c.cache.SetDefault("movie_short:"+strconv.Itoa(m.ID), m)
}

func (c *mediaCache) GetMovieShort(id int) *Movie {
	m, ok := c.cache.Get("movie_short:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return m.(*Movie)
}

func (c *mediaCache) AddTV(t *TVShow) {
	c.cache.SetDefault("tv:"+strconv.Itoa(t.ID), t)
}

func (c *mediaCache) GetTV(id int) *TVShow {
	t, ok := c.cache.Get("tv:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return t.(*TVShow)
}

func (c *mediaCache) AddTVShort(t *TVShow) {
	c.cache.SetDefault("tv_short:"+strconv.Itoa(t.ID), t)
}

func (c *mediaCache) GetTVShort(id int) *TVShow {
	t, ok := c.cache.Get("tv_short:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return t.(*TVShow)
}

func (c *mediaCache) AddEpisode(e *TVEpisode) {
	c.cache.SetDefault("episode:"+strconv.Itoa(e.TVShowID)+":"+strconv.Itoa(e.SeasonNumber)+":"+strconv.Itoa(e.EpisodeNumber), e)
}

func (c *mediaCache) GetEpisode(tvID int, seasonNumber int, episodeNumber int) *TVEpisode {
	e, ok := c.cache.Get("episode:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber) + ":" + strconv.Itoa(episodeNumber))
	if !ok {
		return nil
	}
	return e.(*TVEpisode)
}

func (c *mediaCache) AddSeason(tvID int, seasonNumber int, s []*TVEpisode) {
	c.cache.SetDefault("season:"+strconv.Itoa(tvID)+":"+strconv.Itoa(seasonNumber), s)
}

func (c *mediaCache) GetSeason(tvID int, seasonNumber int) []*TVEpisode {
	s, ok := c.cache.Get("season:" + strconv.Itoa(tvID) + ":" + strconv.Itoa(seasonNumber))
	if !ok {
		return nil
	}
	return s.([]*TVEpisode)
}

func (c *mediaCache) AddMovieSearchResults(query string, page int, results *PaginatedMovieResults) {
	c.cache.SetDefault("movie_search:"+query+":"+strconv.Itoa(page), results)
}

func (c *mediaCache) GetMovieSearchResults(query string, page int) *PaginatedMovieResults {
	r, ok := c.cache.Get("movie_search:" + query + ":" + strconv.Itoa(page))
	if !ok {
		return nil
	}
	return r.(*PaginatedMovieResults)
}

func (c *mediaCache) AddTVSearchResults(query string, page int, results *PaginatedTVShowResults) {
	c.cache.SetDefault("tv_search:"+query+":"+strconv.Itoa(page), results)
}

func (c *mediaCache) GetTVSearchResults(query string, page int) *PaginatedTVShowResults {
	r, ok := c.cache.Get("tv_search:" + query + ":" + strconv.Itoa(page))
	if !ok {
		return nil
	}
	return r.(*PaginatedTVShowResults)
}

func (c *mediaCache) AddMovieGenre(genre *Genre) {
	c.cache.SetDefault("movie_genre:"+strconv.Itoa(genre.ID), genre)
}

func (c *mediaCache) GetMovieGenre(id int) *Genre {
	g, ok := c.cache.Get("movie_genre:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return g.(*Genre)
}

func (c *mediaCache) AddTVGenre(genre *Genre) {
	c.cache.SetDefault("tv_genre:"+strconv.Itoa(genre.ID), genre)
}

func (c *mediaCache) GetTVGenre(id int) *Genre {
	g, ok := c.cache.Get("tv_genre:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return g.(*Genre)
}

func (c *mediaCache) AddActor(actor *Actor) {
	c.cache.SetDefault("actor:"+strconv.Itoa(actor.ID), actor)
}

func (c *mediaCache) GetActor(id int) *Actor {
	a, ok := c.cache.Get("actor:" + strconv.Itoa(id))
	if !ok {
		return nil
	}
	return a.(*Actor)
}
