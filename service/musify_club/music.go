package musify_club

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	model "wmenjoy/music/models"
	"wmenjoy/music/service"
)


type musifyClub struct {
	BaseUrl string
	Crawler service.Crawler
}

var _ service.ISite = (*musifyClub)(nil)
func NewSite(crawler service.Crawler) service.ISite{
	return musifyClub{
		BaseUrl: "https://myzcloud.me/%s",
		Crawler: crawler,
	}
}

func (m musifyClub) AlbumListParser() func(Body io.Reader) (interface{}, error) {
	return func(Body io.Reader) (interface{}, error) {
		return nil, nil
	}
}

func (m musifyClub) AlbumInfoParser() func(body io.Reader) (interface{}, error) {
	return func(body io.Reader) (interface{}, error) {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return nil, err
		}

		urlDiv := doc.Find("meta[itemprop=url]")
		id := ""
		if urlDiv != nil {
			id, _ = urlDiv.Attr("content")
		}

		bodyContent := doc.Find("#bodyContent")

		if bodyContent == nil {
			return nil, nil
		}

		albumInfoDiv := bodyContent.Find(".row.justify-content-center")
		if albumInfoDiv == nil {
			return nil, nil
		}
		baseInfoDiv := albumInfoDiv.Find("div.col-auto")
		image := ""
		albumName := ""
		if baseInfoDiv != nil {
			imgDiv := baseInfoDiv.Find("img")
			if imgDiv != nil {
				albumName, _ = imgDiv.Attr("alt")
				image, _ = imgDiv.Attr("data-src")
			}
		}
		genreInfoDiv := albumInfoDiv.Find(".genre__labels")
		var genres []model.GenreInfo
		if genreInfoDiv != nil {
			genres = make([]model.GenreInfo, 0)
			genreInfoDiv.Find("a").Each(func(i int, selection *goquery.Selection) {
				id, _ := selection.Attr("href")
				genres = append(genres, model.GenreInfo{
					Id : id,
					Name: selection.Get(0).FirstChild.Data[1:],
				})
			})
		}

		artistDiv := albumInfoDiv.Find("ul.icon-list.album-info")
		artists := make([]model.ArtistInfo, 0)
		createDate := ""
		if artistDiv != nil {

			artistDiv.Find("a[itemprop=byArtist]").Each(func(i int, selection *goquery.Selection) {
				id, _ := selection.Attr("href")
				genreA :=  selection.Get(0)
				if genreA != nil {
					artists = append(artists, model.ArtistInfo{
						Id: id,
						Name: selection.Get(0).FirstChild.Data,
					})
				}

			})

			timeDiv := albumInfoDiv.Find("time")
			if timeDiv != nil {
				createDate, _ = timeDiv.Attr("datetime")
			}

		}

		playListDiv := bodyContent.Find("div.playlist.playlist--hover")
		songs := make([]model.MusicInfo,0)
		if playListDiv != nil {
			playListDiv.Find("div.playlist__item").Each(func(i int, selection *goquery.Selection) {

				musicArtist, _:=selection.Attr("data-artist")
				musicName, _ :=selection.Attr("data-name")
				songId, _ := selection.Attr("id")
				downloadUrl := ""
				dataPosition := ""
				playDiv := selection.Find("div.playlist__control")
				if playDiv != nil {
					downloadUrl, _ = playDiv.Attr("data-url")
					dataPosition, _ = playDiv.Attr("data-position")
				}
				songs = append(songs, model.MusicInfo{
					Name: musicName,
					Album: albumName,
					Artist: musicArtist,
					Id: songId[len("playerDiv"):],
					Url: downloadUrl,
					Postion: dataPosition,
				})

			})

		}



		return model.AlbumInfo{
			Id: id,
			Name: albumName,
			Image: image,
			FullName : bodyContent.Find("h1").Get(0).FirstChild.Data,
			Artist: m.parseArtist(artists, bodyContent),
			CreateDate: createDate,
			Genre: genres,
			MusicList: songs,
		}, nil
	}
}


func (m musifyClub) GetUrl(path string) string {
	return fmt.Sprintf(m.BaseUrl, path)
}

func (m musifyClub) parseArtist(extractArtist []model.ArtistInfo, bodyContent *goquery.Selection) []model.ArtistInfo {
	if extractArtist != nil {
		return extractArtist
	}

	breakcrumb := bodyContent.Find(".breadcrumb")
	if breakcrumb == nil {
		return nil
	}
	artists := make([]model.ArtistInfo, 0)
	set := false
	breakcrumb.Find("li").Each(func(index int, selection *goquery.Selection) {
		if set {
			return
		}

		meta := selection.Find("meta")
		if meta == nil {
			return
		}
		val, exist := meta.Attr("content")
		if !exist {
			return
		}
		id, _:= selection.Find("a").Attr("href")
		text := selection.Find("span").Get(0).FirstChild.Data
		if  val == "2" && "Artists" != text && "Исполнители" != text{
			artists = append(artists, model.ArtistInfo{
				Id: id,
				Name: "Various Artists",
			})
			set = true
			return
		} else if val == "3" {
			artists = append(artists,  model.ArtistInfo{
				Id: id,
				Name: text,
			})
			set = true
			return
		}
	})

	return artists
}


