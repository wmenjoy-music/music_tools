package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path"
	etc2 "wmenjoy/music/pkg/etc"
	model2 "wmenjoy/music/pkg/models"
	service2 "wmenjoy/music/pkg/service"
	utils2 "wmenjoy/music/pkg/utils"
	_ "wmenjoy/music/service/musify_club"
)

func Download() error {
	config, err := ParseConfig()
	if err != nil {
		return err
	}
	if config.DownloadDir == "" {
		config.DownloadDir = "./songs"
	}

	if exist, err := utils2.PathExists(config.DownloadDir); !exist || err != nil {
		err = os.MkdirAll(config.DownloadDir, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	logrus.Printf("%+v", config)

	crawler := service2.Crawler{}
	urls := config.Urls
	site := service2.GetSite(urls[0])
	albumList := make([]model2.AlbumInfo, 0)
	for _, url := range urls {
		if site.IsAlbumInfoUrl(url) {
			result, err := crawler.ParsePage(url, site.AlbumInfoParser())
			if err != nil {
				return err
			}
			album := result.(model2.AlbumInfo)
			targetDir := service2.BaseAlbumDownloadDir(config.DownloadDir, album)
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
			for _, albumInfo := range result.([]model2.AlbumInfo) {

				targetDir := service2.BaseAlbumDownloadDir(config.DownloadDir, albumInfo)
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
					saveAlumInfo(targetDir, result.(model2.AlbumInfo))
					albumList = append(albumList, result.(model2.AlbumInfo))
				} else {
					albumList = append(albumList, *album)

				}

			}
		}

	}

	//logrus.Printf("%+v", albumList)

	download := service2.NewDownloader()

	for _, album := range albumList {
		download.PrepareDownload(album, config.DownloadDir)
	}

	download.Wait()

	return nil
}

func saveAlumInfo(dir string, album model2.AlbumInfo) {
	data, err := json.Marshal(album)
	if err != nil {
		return
	}

	_ = os.WriteFile(path.Join(dir, "album.txt"), data, fs.ModePerm)
}

func getAlbumInfoFromDir(dir string) *model2.AlbumInfo {
	if exist, err := utils2.PathExists(path.Join(dir, "album.txt")); !exist || err != nil {
		return nil
	}

	data, err := os.ReadFile(path.Join(dir, "album.txt"))
	if err != nil {
		return nil
	}

	album := &model2.AlbumInfo{}
	err = json.Unmarshal(data, album)
	if err != nil || album.Name == "" {
		return nil
	}
	return album
}

func ParseConfig() (*etc2.Config, error) {
	config := &etc2.Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
