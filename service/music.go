package service

import (
	"io"
	model "wmenjoy/music/models"
)

type ISite interface {
	// GetUrl 根据Path 获取绝对的url
	GetUrl(path string) string
	AlbumListParser() func(body io.Reader) (interface{}, error)
	AlbumInfoParser() func(body io.Reader) (interface{}, error)
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
	GetAlbum(id string) ([]model.MusicInfo, error)
}

type DefaultDetailer struct {
}
