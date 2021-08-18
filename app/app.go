package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path"
	"wmenjoy/music/pkg/etc"
	"wmenjoy/music/pkg/models"
	"wmenjoy/music/pkg/service"
	"wmenjoy/music/pkg/utils"
	_ "wmenjoy/music/pkg/service/musify_club"
)

func Download() error {
	config, err := ParseConfig()
	if err != nil {
		return err
	}
	if config.DownloadDir == "" {
		config.DownloadDir = "./songs"
	}

	if exist, err := utils.PathExists(config.DownloadDir); !exist || err != nil {
		err = os.MkdirAll(config.DownloadDir, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	logrus.Printf("%+v", config)

	crawler := service.Crawler{}
	urls := config.Urls
	site := service.GetSite(urls[0])
	albumList := make([]model.AlbumInfo, 0)
	for _, url := range urls {
		if site.IsAlbumInfoUrl(url) {
			result, err := crawler.ParsePage(url, site.AlbumInfoParser())
			if err != nil {
				return err
			}
			album := result.(model.AlbumInfo)
			targetDir := service.BaseAlbumDownloadDir(config.DownloadDir, album)
			err = os.MkdirAll(targetDir, fs.ModePerm)
			if err != nil {
				return err
			}
			saveAlumInfo(targetDir, album)
			albumList = append(albumList, album)
		} else {
			result, err := crawler.ParsePage(site.NormalUrl(url), site.AlbumListParser())
			if err != nil {
				return err
			}
			for _, albumInfo := range result.([]model.AlbumInfo) {

				targetDir := service.BaseAlbumDownloadDir(config.DownloadDir, albumInfo)
				err = os.MkdirAll(targetDir, fs.ModePerm)
				if err != nil {
					return err
				}
				album := getAlbumInfoFromDir(targetDir)

				if album == nil {
					result, err = crawler.ParsePage(albumInfo.Url, site.AlbumInfoParser())
					if err != nil {
						return err
					}
					saveAlumInfo(targetDir, result.(model.AlbumInfo))
					albumList = append(albumList, result.(model.AlbumInfo))
				} else {
					albumList = append(albumList, *album)

				}

			}
		}

	}

	//logrus.Printf("%+v", albumList)

	download := service.NewDownloader()

	for _, album := range albumList {
		download.PrepareDownload(album, config.DownloadDir)
	}

	download.CloseDataChannel()

	download.Wait()

	return nil
}

func saveAlumInfo(dir string, album model.AlbumInfo) {
	data, err := json.Marshal(album)
	if err != nil {
		return
	}

	_ = os.WriteFile(path.Join(dir, "album.txt"), data, fs.ModePerm)
}

func getAlbumInfoFromDir(dir string) *model.AlbumInfo {
	if exist, err := utils.PathExists(path.Join(dir, "album.txt")); !exist || err != nil {
		return nil
	}

	data, err := os.ReadFile(path.Join(dir, "album.txt"))
	if err != nil {
		return nil
	}

	album := &model.AlbumInfo{}
	err = json.Unmarshal(data, album)
	if err != nil || album.Name == "" {
		return nil
	}
	return album
}

func ParseConfig() (*etc.Config, error) {
	config := &etc.Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
