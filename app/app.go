package app

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"wmenjoy/music/etc"
	model "wmenjoy/music/models"
	"wmenjoy/music/service"
)

func Download() error{
	config, err:= ParseConfig()
	if err != nil {
		return err
	}

	logrus.Printf("%+v", config)

	crawler := service.Crawler{}
	urls :=config.Urls
	site := service.GetSite(urls[0])
    albumList := make([]model.AlbumInfo, 0)
	for _, url := range urls {
		if site.IsAlbumInfoUrl(url){
			result, err := crawler.ParsePage(url, site.AlbumInfoParser())
			if err != nil {
				return err
			}
			albumList = append(albumList, result.(model.AlbumInfo))
		} else {
			result, err := crawler.ParsePage(site.NormalUrl(url), site.AlbumListParser())
			if err != nil {
				return err
			}
			for _, albumInfo := range result.([]model.AlbumInfo) {

				result, err := crawler.ParsePage(albumInfo.Url, site.AlbumInfoParser())
				if err != nil {
					return err
				}
				albumList = append(albumList, result.(model.AlbumInfo))
			}
		}
		
	}

	logrus.Printf("%+v", albumList)
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
