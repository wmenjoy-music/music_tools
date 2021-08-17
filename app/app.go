package app

import (
	"strings"
	"wmenjoy/music/etc"
	model "wmenjoy/music/models"
	"wmenjoy/music/service"
	"wmenjoy/music/service/musify_club"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Download() error{
	config, err:= ParseConfig()
	if err != nil {
		return err
	}

	logrus.Printf("%+v", config)

	crawler := service.Crawler{}
	site := musify_club.NewSite()
	urls :=config.Urls
    albumList := make([]model.AlbumInfo, 0)
	for _, url := range urls {
		if strings.HasPrefix(url, "https://w1.musify.club/release") ||
		  strings.HasPrefix(url, "https://w1.musify.club/en/release"){
			result, err := crawler.ParsePage(url, site.AlbumInfoParser())
			if err != nil {
				return err
			}
			albumList = append(albumList, result.(model.AlbumInfo))
		} else {
			realUrl := url
			if !strings.HasSuffix(realUrl, "/release")}{
				realUrl = realUrl + "/release"
			} 

			result, err := crawler.ParsePage(url, site.AlbumInfoParser())
			if err != nil {
				return err
			}
		
			albumList = append(albumList, result.([]model.MusicInfo))
			musicList := crawler.ParsePage(url, )
		}
		
	}

	//service.PrepareDownload(,)
	

	return nil
}

func ParseConfig() (*etc.Config, error) {
	config := &etc.Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
