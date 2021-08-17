package musify_club

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	model "wmenjoy/music/models"
)

func TestMusifyClub_AlbumInfoParser(t *testing.T) {
	musifyClub := NewSite()

	//crawler := service.Crawler{}
	file, err := os.Open("/Users/liujinliang/workspace/music/music_manager/test.html")
	assert.Nil(t, err)
	p, err := musifyClub.AlbumInfoParser()(file)

	//p, err := crawler.ParsePage("https://w1.musify.club/release/medlyaki-shkolnih-diskotek-2021-1484710", )

	assert.Nil(t, err)
	assert.Equal(t, "Медляки Школьных Дискотек (2021)",p.(model.AlbumInfo).FullName)
}
