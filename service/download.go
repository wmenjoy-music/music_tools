package service

import (
	"fmt"
	"os"
	"sync"
	model "wmenjoy/music/models"
)

var songChan = make(chan downloadInfo, 10)

var closeChan = make(chan struct{}, 10)

var sigChan = make(chan os.Signal)

var crawler Crawler

type downloadInfo struct {
	object IDownloadObject
	downloadDir  string
}

func Start(threadNum int){
	wg := sync.WaitGroup{}
	wg.Add(threadNum)
	for i :=0; i < threadNum; i++{
		go Run(wg)
	}
	wg.Wait()
}


func Close(){
	close(songChan)
	close(closeChan)
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
		case x, ok := <-songChan:
			if !ok {
				group.Done()
				return
			}
			err := crawler.Download(x.object, x.downloadDir)
			if err != nil {
				//
			}

		case _, ok := <-closeChan:
			if !ok {
				group.Done()
			}
			return
		case sig := <-sigChan:
			fmt.Println("接受到来自系统的信号：",sig.String())
			Close()
		}
	}
}