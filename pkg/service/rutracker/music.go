package rutracker

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Rutracker struct {
	baseUrl string
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
