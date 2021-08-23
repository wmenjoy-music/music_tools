package rutracker

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"wmenjoy/music/pkg/model"
	"wmenjoy/music/pkg/service"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Rutracker struct {
	baseUrl string
}

func NewForumSite() service.IForumSite {
	return Rutracker{
		baseUrl: "https://rutracker.org/%s",
	}
}

// SubForumParser 重要的是要解析什么？
func (r Rutracker) SubForumParser() func(body io.Reader) (interface{}, error) {
	return func(body io.Reader) (interface{}, error) {

		return nil, nil
	}
}

func (r Rutracker) GetUrl(path string) string {
	if strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "http://") {
		return path
	}

	return fmt.Sprintf(r.baseUrl, path)

}

var _ service.IForumSite = (*Rutracker)(nil)

func (r Rutracker) DetailParser() func(body io.Reader) (interface{}, error) {
	return func(body io.Reader) (interface{}, error) {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(transform.NewReader(body, charmap.Windows1251.NewDecoder()))
		if err != nil {
			return nil, err
		}

		socDiv := doc.Find("#soc-container")

		dataShareUrl, _ := socDiv.Attr("data-share_url")
		dataShareTitle, _ := socDiv.Attr("data-share_title")

		divMessage := doc.Find("#topic_main tbody[id]").First().Find("td.message")

		if divMessage == nil {
			return nil, errors.New("没有信息")
		}

		postWrapDiv := divMessage.Find("div.post_wrap").First()
		if postWrapDiv == nil {
			return nil, errors.New("没有信息")
		}

		postBody := postWrapDiv.Find("div>.post_body")
		extraLinkData, _ := postBody.Attr("data-ext_link_data")

		text := postBody.Text()

		aMagnet := divMessage.Find("a.magnet-link")
		magnetTitle, _ := aMagnet.Attr("title")
		magnetLink, _ := aMagnet.Attr("href")
		aTorrent := divMessage.Find("a.dl-stub.dl-link.dl-topic")
		torrentUrl, _ := aTorrent.Attr("href")
		logrus.Printf("%s,%s,%s,%s", dataShareUrl, dataShareTitle, extraLinkData, text)
		tags := []string{"Жанр", "Страна исполнителя (группы)", "Год издания", "Аудиокодек", "Тип рипа", "Битрейт аудио", "Продолжительность", "Исполнитель",
			"Альбом", "Страна", "Дата выпуска", "Формат", "Битрейт", "Издатель (лейбл)", "Носитель", "Номер по каталогу", "Годы выпуска дисков", "Лейбл/Label",
		}
		tagMap := r.parseTag(text, tags)

		getTag := func(tagMap map[string]string, keys ...string) string {
			for _, key := range keys {
				if tagMap[key] != "" {
					return tagMap[key]
				}
			}
			return ""
		}

		if strings.Contains(dataShareTitle, "Дискография") ||
			strings.Contains(dataShareTitle, "дискография") ||
			strings.Contains(dataShareTitle, "Discography") ||
			strings.Contains(dataShareTitle, "discography") ||
			strings.Contains(dataShareTitle, "Collection") ||
			strings.Contains(dataShareTitle, "Коллекция") {
			if (strings.Contains(dataShareTitle, "Collection") ||
				strings.Contains(dataShareTitle, "Коллекция")) &&
				strings.Contains(dataShareTitle, "VA") {
				return nil, nil
			}
			return &model.ForumDiscographyInfo{
				Artist:      getTag(tagMap, "Исполнитель", "title"),
				Title:       dataShareTitle,
				Url:         dataShareUrl,
				Years:       getTag(tagMap, "Год выпуска диска", "Годы выпуска дисков", "Год издания"),
				GenreType:   getTag(tagMap, "Жанр"),
				Country:     getTag(tagMap, "Страна исполнителя (группы)", "Страна"),
				BitRate:     getTag(tagMap, "Битрейт аудио", "Битрейт"),
				FileType:    getTag(tagMap, "Аудиокодек", "Формат"),
				Duration:    tagMap["Продолжительность"],
				Content:     text,
				MagnetLink:  magnetLink,
				MagnetTitle: magnetTitle,
				Torrent:     r.GetUrl(torrentUrl),
				ExtraInfo:   r.ParseAlbumInfo(postBody),
			}, nil

		}

		return &model.ForumAlbumInfo{
			Artist:      getTag(tagMap, "Исполнитель", "title"),
			Name:        getTag(tagMap, "Альбом", "title"),
			Title:       dataShareTitle,
			Url:         dataShareUrl,
			Year:        getTag(tagMap, "Дата выпуска"),
			GenreType:   getTag(tagMap, "Жанр"),
			Country:     getTag(tagMap, "Страна исполнителя (группы)", "Страна"),
			BitRate:     getTag(tagMap, "Битрейт аудио", "Битрейт"),
			FileType:    getTag(tagMap, "Аудиокодек", "Формат"),
			Duration:    tagMap["Продолжительность"],
			Content:     text,
			MagnetLink:  magnetLink,
			MagnetTitle: magnetTitle,
			Torrent:     r.GetUrl(torrentUrl),
			ExtraInfo:   r.ParseAlbumInfo(postBody),
		}, nil

	}
}

func (r Rutracker) ItemListParser() func(body io.Reader) (interface{}, error) {
	return func(body io.Reader) (interface{}, error) {

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(transform.NewReader(body, charmap.Windows1251.NewDecoder()))
		if err != nil {
			return nil, err
		}
		divContent := doc.Find("#main_content_wrap")
		if divContent == nil {
			return nil, errors.New("不存在main_content_wrap")
		}

		divContent.Find("table.vf-table.vf-tor.forumline.forum tr").Each(func(i int, s *goquery.Selection) {
			topicId, _ := s.Attr("data-topic_id")
			id, _ := s.Attr("id")
			if id != "" {
				logrus.Printf("%s:%s:%s", topicId, id, s.Find("td.vf-col-t-title div.torTopic a.tt-text").Text())
			}

		})

		return nil, nil
	}
}

func (r Rutracker) parseTag(text string, tags []string) map[string]string {
	tagMap := make(map[string]string, 0)

	lines := strings.Split(text, "\n")

	parseTag := func(line string) {
		for _, tag := range tags {
			if strings.Contains(line, tag) {
				value := strings.TrimSpace(line[strings.Index(line, ":")+1:])

				if strings.HasSuffix(value, "Треклист:") {
					value = value[:len(value)-len("Треклист:")]
					tagMap["type"] = "album"
				} else if strings.HasSuffix(value, "Tracklist:") {
					value = value[:len(value)-len("Tracklist:")]
					tagMap["type"] = "album"
				} else if strings.HasSuffix(value, "Albums:") {
					value = value[:len(value)-len("Albums:")]
					tagMap["type"] = "artist"
				}
				tagMap[tag] = value
			}
		}

	}

	for _, line := range lines {
		if strings.Contains(line, "|") && strings.Contains(line, ":") {
			subTags := strings.Split(line, "|")
			for _, subTag := range subTags {
				parseTag(subTag)
			}
		} else if index := strings.Index(line, "- Дискография"); index >= 0 {
			tagMap["title"] = strings.TrimSpace(line[:index])
		} else if index := strings.Index(line, "- дискография"); index >= 0 {
			tagMap["title"] = strings.TrimSpace(line[:index])
		} else if index := strings.Index(line, "- Discography"); index >= 0 {
			tagMap["title"] = strings.TrimSpace(line[:index])
		} else if index := strings.Index(line, "- discography"); index >= 0 {
			tagMap["title"] = strings.TrimSpace(line[:index])
		} else {
			parseTag(line)
		}

	}
	return tagMap
}

func (r Rutracker) ParseAlbumInfo(postBody *goquery.Selection) map[string]interface{} {
	return r.parseSpWrapper(postBody.ChildrenFiltered(".sp-wrap"))
}

func (r Rutracker) parseSpWrapper(selection *goquery.Selection) map[string]interface{} {
	totalmap := make(map[string]interface{})
	var list []interface{} = nil
	listTitle := ""
	selection.Each(func(i int, selection *goquery.Selection) {
		// 如果之前有span
		prevSelection := selection.Prev()
		if prevSelection.Size() > 0 && (prevSelection.Find("span.post-i").Size() > 0 || prevSelection.Find("span.p-color").Size() > 0) && prevSelection.Get(0).Data == "span" && strings.Contains(prevSelection.Text(), ":") {

			if list != nil {
				totalmap[listTitle] = list
			}
			list = make([]interface{}, 0)
			listTitle = strings.TrimSpace(prevSelection.Text())
			listTitle = listTitle[:len(listTitle)-1]
		}

		if selection.Find(".sp-body .sp-wrap").Size() > 0 {
			title := selection.Find("div.sp-head").First().Text()
			subMap := r.parseSpWrapper(selection.Find(".sp-body").First().ChildrenFiltered(".sp-wrap"))
			if list != nil {
				list = append(list, subMap)
			} else {
				totalmap[title] = subMap
			}
			return
		}

		title := selection.Find("div.sp-head").First().Text()
		sp_body := selection.Find("div.sp-body")
		sp_text := sp_body.Text()
		image, _ := sp_body.Find(".postImg").Attr("title")
		detail := model.ForumAlbumDetail{
			Name:    title,
			Content: sp_text,
			Image:   image,
		}
		if list != nil {
			list = append(list, detail)
		} else {
			totalmap[title] = detail
		}
	})

	if list != nil {
		totalmap[listTitle] = list
	}

	return totalmap
}
