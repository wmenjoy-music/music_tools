package service

import (
	"sync"
	model "wmenjoy/music/models"
)

var songChan = make(chan downloadInfo, 10)

var signalChan = make(chan struct{}, 10)

var crawler Crawler

type downloadInfo struct {
	object IDownloadObject
	downloadDir  string
}

func Close(){
	close(songChan)
	signalChan <- struct{}{}
}

// PrepareDownload 准备目录， 将下载数据发送到Channels
func PrepareDownload(info model.AlbumInfo, baseDir string) {

	songChan <- downloadInfo{
		object: DownloadImage{
			DownloadUrl: info.Image,
			FileType: "jpg",
			Name: "cover",
		},
	}

	for _, song := range info.MusicList{
		songChan <- downloadInfo{
			object: DownloadMusic{
				DownloadUrl: song.Url,
				FileType: "mp3",
				Name: song.Name,
				Artist: song.Artist,
				index: song.Postion,
				Category: info.Category,
			},
		}
	}
}

func Run(group sync.WaitGroup){
	for {
		select {
		case x := <-songChan:
			err := crawler.Download(x.object, x.downloadDir)
			if err != nil {
				//
			}

		case <-signalChan:
			group.Done()
			return
		}
	}
}