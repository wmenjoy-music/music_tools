package model

type ArtistInfo struct {
	Name    string
	Memeber []string
	Country string
}

type MusicInfo struct {
	Name        string
	Url         string
	Album       string
	Artist      string
	Postion     string
	DataTitle   string
	DownloadUrl string
	Download    string
}

type AlbumInfo struct {
	Name       string
	FullName   string
	Image      string
	ArtistName []string
	Genre      string
	musicList  []string
	dataType   string
	year       string
	category   string
	url        string
}
