package musify_club

import (
	"fmt"
	"io"
	"strings"
	model "wmenjoy/music/models"
	"wmenjoy/music/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type musifyClub struct {
	BaseUrl string
}

var _ service.ISite = (*musifyClub)(nil)

func NewSite() service.ISite {
	return musifyClub{
		BaseUrl: "https://myzcloud.me/%s",
	}
}

func (m musifyClub) AlbumListParser() func(Body io.Reader) (interface{}, error) {
	return func(body io.Reader) (interface{}, error) {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return nil, err
		}

		bodyContent := doc.Find("#bodyContent")
		if bodyContent == nil {
			return nil, nil
		}
		divAlbumList := bodyContent.Find("div#divAlbumsList")
		if divAlbumList == nil {
			return nil, nil
		}
		albums := make([]model.AlbumInfo, 0)

		cards := divAlbumList.Find(".card.release-thumbnail")
		if cards == nil {
			return nil, nil
		}
		cards.Each(func(index int, selection *goquery.Selection) {
			dataType, _ := selection.Attr("data-type")
			urlAddrA := selection.Find("a")
			image := ""
			url := ""
			name := ""
			id := ""
			year := ""
			genres := make([]model.GenreInfo, 0)

			if urlAddrA != nil {
				id, _ = urlAddrA.Attr("href")
				url = m.GetUrl(id)
				img := urlAddrA.Find("img")
				if img != nil {
					image, _ = img.Attr("data-src")
				}
				name, _ = img.Attr("alt")
				if name == "垃圾" {
					return
				}

			}
			yeara := selection.Find("p.card-text a")
			if yeara != nil {
				year = yeara.Get(0).FirstChild.Data
			}

			pCardGenre := selection.Find("p.card-text.genre__labels")

			if pCardGenre != nil {
				pCardGenre.Find("a").Each(func(i int, s *goquery.Selection) {
					gid, _ := s.Attr("href")
					genres = append(genres, model.GenreInfo{
						Id:   gid,
						Name: s.Get(0).FirstChild.Data,
						Url:  m.GetUrl(gid),
					})
				})
			}

			iCalendar := selection.Find("i.zmdi.zmdi-calendar")
			createDate := ""
			if iCalendar != nil {
				createDate = strings.TrimSpace(iCalendar.Parent().Text())
			}
			iRating := selection.Find("i.zmdi.zmdi-star.zmdi-hc-fw")
			rating := ""
			if iRating != nil {
				rating = strings.TrimSpace(iRating.Parent().Text())
			}

			albums = append(albums, model.AlbumInfo{
				Name:       name,
				Id:         id,
				Url:        url,
				Genre:      genres,
				Image:      image,
				DataType:   dataType,
				Artist:     []model.ArtistInfo{*m.parseAlbumListArtist(bodyContent)},
				Year:       year,
				CreateDate: createDate,
				Rating:     rating,
				Category:   model.CategoryTypeMap[dataType],
			})
		})

		return albums, nil
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
					Id:   id,
					Name: selection.Get(0).FirstChild.Data[1:],
					Url:  m.GetUrl(id),
				})
			})
		}

		artistDiv := albumInfoDiv.Find("ul.icon-list.album-info")
		artists := make([]model.ArtistInfo, 0)
		createDate := ""
		if artistDiv != nil {

			artistDiv.Find("a[itemprop=byArtist]").Each(func(i int, selection *goquery.Selection) {
				id, _ := selection.Attr("href")
				genreA := selection.Get(0)
				if genreA != nil {
					artists = append(artists, model.ArtistInfo{
						Id:   id,
						Name: selection.Get(0).FirstChild.Data,
						Url:  m.GetUrl(id),
					})
				}

			})

			timeDiv := albumInfoDiv.Find("time")
			if timeDiv != nil {
				createDate, _ = timeDiv.Attr("datetime")
			}

		}

		playListDiv := bodyContent.Find("div.playlist.playlist--hover")
		songs := make([]model.MusicInfo, 0)
		if playListDiv != nil {
			playListDiv.Find("div.playlist__item").Each(func(i int, selection *goquery.Selection) {

				musicArtist, _ := selection.Attr("data-artist")
				musicName, _ := selection.Attr("data-name")
				songId, _ := selection.Attr("id")
				downloadUrl := ""
				dataPosition := ""
				playDiv := selection.Find("div.playlist__control")
				if playDiv != nil {
					downloadUrl, _ = playDiv.Attr("data-url")
					dataPosition, _ = playDiv.Attr("data-position")
				}
				spanDelete := selection.Find("span.badge.badge-pill.badge-danger")
				if spanDelete != nil {
					logrus.Print(musicName + " is " + spanDelete.Get(0).FirstChild.Data)
					return
				}

				songs = append(songs, model.MusicInfo{
					Name:    musicName,
					Album:   albumName,
					Artist:  musicArtist,
					Id:      songId[len("playerDiv"):],
					Url:     downloadUrl,
					Postion: dataPosition,
				})

			})

		}

		return model.AlbumInfo{
			Id:         id,
			Name:       albumName,
			Image:      image,
			FullName:   bodyContent.Find("h1").Get(0).FirstChild.Data,
			Artist:     m.parseAlbumInfoArtist(artists, bodyContent),
			CreateDate: createDate,
			Genre:      genres,
			MusicList:  songs,
		}, nil
	}
}

func (m musifyClub) GetUrl(path string) string {
	return fmt.Sprintf(m.BaseUrl, path)
}

func (m musifyClub) parseAlbumListArtist(bodyContent *goquery.Selection) *model.ArtistInfo {

	breakcrumb := bodyContent.Find(".breadcrumb")
	if breakcrumb == nil {
		return nil
	}
	var artist *model.ArtistInfo
	breakcrumb.Find("li").Each(func(index int, selection *goquery.Selection) {
		meta := selection.Find("meta")
		if meta == nil {
			return
		}
		val, exist := meta.Attr("content")
		if !exist {
			return
		}
		id, _ := selection.Find("a").Attr("href")
		text := selection.Find("span").Get(0).FirstChild.Data
		if val == "3" {
			artist = &model.ArtistInfo{
				Id:   id,
				Name: text,
				Url:  m.GetUrl(id),
			}
		}
	})

	return artist
}

func (m musifyClub) parseAlbumInfoArtist(extractArtist []model.ArtistInfo, bodyContent *goquery.Selection) []model.ArtistInfo {
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
		id, _ := selection.Find("a").Attr("href")
		text := selection.Find("span").Get(0).FirstChild.Data
		if val == "2" && "Artists" != text && "Исполнители" != text {
			artists = append(artists, model.ArtistInfo{
				Id:   id,
				Name: "Various Artists",
			})
			set = true
			return
		} else if val == "3" {
			artists = append(artists, model.ArtistInfo{
				Id:   id,
				Name: text,
				Url:  m.GetUrl(id),
			})
			set = true
			return
		}
	})

	return artists
}
