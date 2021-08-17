package musify_club

import (
	"os"
	"testing"
	model "wmenjoy/music/models"

	"github.com/stretchr/testify/assert"
)

func TestMusifyClub_AlbumInfoParser(t *testing.T) {
	musifyClub := NewSite()

	//crawler := service.Crawler{}
	file, err := os.Open("/Users/liujinliang/workspace/music/music_manager/test/test.html")
	assert.Nil(t, err)
	p, err := musifyClub.AlbumInfoParser()(file)

	//p, err := crawler.ParsePage("https://w1.musify.club/release/medlyaki-shkolnih-diskotek-2021-1484710", )

	assert.Nil(t, err)
	assert.Equal(t, "Медляки Школьных Дискотек (2021)", p.(model.AlbumInfo).FullName)
}

func TestMusifyClub_AlbumListParser(t *testing.T) {
	musifyClub := NewSite()
	file, err := os.Open("d:/code/music_tools/test/albums_list.html")
	assert.Nil(t, err)
	p, err := musifyClub.AlbumListParser()(file)
	assert.Nil(t, err)
	assert.Equal(t, "Медляки Школьных Дискотек (2021)", p.(model.AlbumInfo).FullName)
}
