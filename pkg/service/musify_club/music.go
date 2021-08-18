package musify_club

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	model "wmenjoy/music/pkg/models"
	service "wmenjoy/music/pkg/service"
)

type musifyClub struct {
	BaseUrl string
}

func init() {
	service.RegistSite("musify_club", NewSite)
}

var _ service.ISite = (*musifyClub)(nil)

func NewSite() service.ISite {
	return musifyClub{
		BaseUrl: "https://w1.musify.club%s",
	}
}

func (m musifyClub) NormalUrl(url string) string {
	if !strings.HasSuffix(url, "/releases") {
		return url + "/releases"
	}
	return url
}
func (m musifyClub) IsAlbumInfoUrl(url string) bool {
	return strings.HasPrefix(url, "https://w1.musify.club/release") ||
		strings.HasPrefix(url, "https://w1.musify.club/en/release")
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
		typeI := albumInfoDiv.Find("i.zmdi.zmdi-collection-music.zmdi-hc-fw")
		category := ""
		if typeI != nil {
			category = strings.TrimSpace(typeI.Parent().Text())
		}

		artistDiv := albumInfoDiv.Find("ul.icon-list.album-info")
		artists := make([]model.ArtistInfo, 0)
		createDate := ""
		year := ""
		if artistDiv != nil {

			artistDiv.Find("a[itemprop=byArtist]").Each(func(i int, selection *goquery.Selection) {
				id, _ := selection.Attr("href")

				artists = append(artists, model.ArtistInfo{
					Id:   id,
					Name: strings.TrimSpace(selection.Text()),
					Url:  m.GetUrl(id),
				})

			})

			timeDiv := albumInfoDiv.Find("time")

			if timeDiv != nil {
				createDate, _ = timeDiv.Attr("datetime")
				year = timeDiv.Parent().Find("a").Text()
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
				songId = songId[len("playerDiv"):]
				playDiv := selection.Find("div#play_" + songId)
				dataTitle := ""
				if playDiv != nil {
					downloadUrl, _ = playDiv.Attr("data-url")
					dataPosition, _ = playDiv.Attr("data-position")
					dataTitle, _ = playDiv.Attr("data-title")
				}

				songUrl, _ := selection.Find("div.playlist__heading a.strong").Attr("href")
				stars := selection.Find("i.zmdi.zmdi-star-circle").Parent().Text()

				trackDetailDiv := selection.Find("div.track__details")
				duration := ""
				bitrate := ""
				if trackDetailDiv != nil {
					trackDetailDiv.Find("span").Each(func(i int, selection *goquery.Selection) {
						if strings.Contains(selection.Text(), "Кб/с") {
							bitrate = strings.TrimSpace(strings.ReplaceAll(selection.Text(), "Кб/с", "Kbps"))
						} else if strings.Contains(selection.Text(), ":") {
							duration = selection.Text()
						}
					})

				}

				spanDelete := selection.Find("span.badge.badge-pill.badge-danger")
				if downloadUrl == "" || (spanDelete != nil && spanDelete.Size() != 0 && spanDelete.Get(0).FirstChild.Data == "Недоступен") {
					logrus.Print(musicName + " is deleted")
					return
				}

				songs = append(songs, model.MusicInfo{
					Name:        musicName,
					Album:       albumName,
					Artist:      musicArtist,
					Id:          songId,
					Url:         m.GetUrl(songUrl),
					Postion:     dataPosition,
					DownloadUrl: m.GetUrl(downloadUrl),
					DataTitle:   dataTitle,
					Stars:       stars,
					BitRate:     bitrate,
					Duration:    duration,
				})

			})

		}

		rating, _ := bodyContent.Find("select#rating").Attr("data-rating")

		return model.AlbumInfo{
			Id:         id,
			Url:        m.GetUrl(id),
			Name:       albumName,
			Image:      image,
			FullName:   bodyContent.Find("h1").Get(0).FirstChild.Data,
			Artist:     m.parseAlbumInfoArtist(artists, bodyContent),
			CreateDate: createDate,
			Year:       year,
			Genre:      genres,
			MusicList:  songs,
			Category:   model.NormalCategory(category),
			DataType:   model.TypeCategoryMap[category],
			Rating:     rating,
		}, nil
	}
}

func (m musifyClub) GetUrl(path string) string {
	if strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "http://") {
		return path
	}
	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf(m.BaseUrl, path)
	} else {
		return fmt.Sprintf(m.BaseUrl+"/", path)
	}

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
	if extractArtist != nil && len(extractArtist) > 0 {
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
