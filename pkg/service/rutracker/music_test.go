package rutracker

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRutracker_ParserAlum(t *testing.T) {
	r := NewForumSite()

	//crawler := service.Crawler{}
	file, err := os.Open("d://code/music_tools/test/rutracker_forum_list.html")
	assert.Nil(t, err)
	p, err := r.ItemListParser()(file)

	//p, err := crawler.ParsePage("https://w1.musify.club/release/medlyaki-shkolnih-diskotek-2021-1484710", )

	assert.Nil(t, err)
	assert.Nil(t, p)
}
func TestRutracker_ParserForumPost(t *testing.T) {
	r := NewForumSite()

	//crawler := service.Crawler{}
	file, err := os.Open("d://code/music_tools/test/rutracker_forum_detail_2.html")
	assert.Nil(t, err)
	p, err := r.DetailParser()(file)

	//p, err := crawler.ParsePage("https://w1.musify.club/release/medlyaki-shkolnih-diskotek-2021-1484710", )

	assert.Nil(t, err)
	assert.NotNil(t, p)
}
