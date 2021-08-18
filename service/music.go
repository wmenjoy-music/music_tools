package service

import (
	"io"
	"strings"
	model "wmenjoy/music/models"
)

type ISite interface {
	IsAlbumInfoUrl(url string) bool
	// GetUrl 根据Path 获取绝对的url
	GetUrl(path string) string
	AlbumListParser() func(body io.Reader) (interface{}, error)
	AlbumInfoParser() func(body io.Reader) (interface{}, error)
	NormalUrl(url string) string
}

type SiteFactory func() ISite

var siteRegistry map[string]SiteFactory = make(map[string]SiteFactory, 1)

func RegistSite(name string, factory SiteFactory) {
	siteRegistry[name] = factory
}

func GetSite(url string) ISite {
	if strings.Contains(url, "w1.musify.club") {
		return siteRegistry["musify_club"]()
	}
	return nil
}

type ListOptions struct {
	Start int
	End   int
}

type Lister interface {
	// GetGenres 获取所有
	GetGenres(options ...ListOptions) (model.GenreInfo, error)
	GetAlbumsByGenre(id string, options ...ListOptions) ([]model.AlbumInfo, error)
	GetAlbumsByYear(year string, options ...ListOptions) ([]model.AlbumInfo, error)
	GetAlbumsByArtist(id string, options ...ListOptions) ([]model.AlbumInfo, error)
	GetArtists(options ...ListOptions) (model.ArtistInfo, error)
}

type Detailer interface {
	GetAlbum(id string) (model.AlbumInfo, error)
}

type DefaultDetailer struct {
}
