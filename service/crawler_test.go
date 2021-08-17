package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"io"
	"io/fs"
	"os"
	"testing"
)

func TestCrawler_DownloadObj(t *testing.T) {
	crawler := &Crawler{
		Options: Options{
			ShowProgress: true,
		},
	}

	assert.NotNil(t, crawler)

	os.MkdirAll("/Users/liujinliang/workspace/music/music_manager/tmp/Live At Knebworth", fs.ModePerm)

	err := crawler.Download(DownloadMusic{
		Artist: "Pink Floyd",
		Name: "Shine On You Crazy Diamond (Parts 1-5) (Live 1990)",
		FileType: "mp3",
		index: "1",
		DownloadUrl: "https://w1.musify.club/track/dl/16133463/pink-floyd-shine-on-you-crazy-diamond-parts-1-5-live-1990.mp3",

	}, "/Users/liujinliang/workspace/music/music_manager/tmp/Live At Knebworth")
	assert.Nil(t, err)

}

func TestCrawler_ParsePage(t *testing.T) {
	crawler := &Crawler{
		Options: Options{
			ShowProgress: true,
		},
	}
	assert.NotNil(t, crawler)

	b , err := crawler.ParsePage("http://metalsucks.net", func(body io.Reader) (interface{}, error) {
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return nil ,err
		}
		b := ""
		// Find the review items
		doc.Find(".left-content article .post-title").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the title
			title := s.Find("a").Text()
			b = title
			fmt.Printf("Review %d: %s\n", i, title)
		})
		return b, nil
	})
	assert.Nil(t, err)
	assert.NotNil(t, b)

}
