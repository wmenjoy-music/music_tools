package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vbauerster/mpb/v7"
	"io/fs"
	"os"
	"path"
	"sync"
	"wmenjoy/music/pkg/model"
	"wmenjoy/music/pkg/utils"
)

type Downloader struct {
	songChan  chan downloadInfo
	closeChan chan struct{}
	sigChan   chan os.Signal
	crawler   Crawler
	wg        *sync.WaitGroup
	bars      []*mpb.Bar
	barWaitGroup  []*sync.WaitGroup
	doneWg    *sync.WaitGroup
}

func NewDownloader() *Downloader {
	channelLength := viper.GetInt("songChannelLength")

	if channelLength <= 0 {
		channelLength = 1000
	}
	retry := viper.GetInt("retryCount")
	if retry <= 0 {
		retry = 3
	}

	threadNum := viper.GetInt("threadNum")
	doneWg := new(sync.WaitGroup)
	barWaitGroup := make([]*sync.WaitGroup, threadNum)
	for i := 0; i < threadNum; i++{
		barWaitGroup[i] = new(sync.WaitGroup)
	}

	download := &Downloader{
		songChan:  make(chan downloadInfo, channelLength),
		closeChan: make(chan struct{}, 1),
		sigChan:   make(chan os.Signal),
		crawler: Crawler{
			Retry: retry,
			Options: Options{
				ShowProgress: true,
			},
			ProcessBars: mpb.New(),
		},
		wg: &sync.WaitGroup{},
		doneWg: doneWg,
		bars: make([]*mpb.Bar, threadNum),
		barWaitGroup: barWaitGroup,
	}
	//d.Start(viper.GetInt("threadNum"))
	download.Start(threadNum)
	return download
}

type downloadInfo struct {
	object      IDownloadObject
	downloadDir string
}

func (d *Downloader) Start(threadNum int) {
	logrus.Printf("threadNum is %d", threadNum)
	d.wg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go d.Run(d.wg)
	}


}

func (d *Downloader) CloseDataChannel() {
	close(d.songChan)
}

func (d *Downloader) Wait() {

	d.wg.Wait()
}

func (d* Downloader) Close() {
	close(d.songChan)
	close(d.closeChan)
}

func BaseAlbumDownloadDir(baseDir string, info model.AlbumInfo) string {

	artist := "VA"

	if len(info.Artist) == 1 && info.Artist[0].Name != "Various Artists" {
		artist = utils.ValidateFileName(info.Artist[0].Name)
	}

	dirName := utils.ValidateFileName(info.Name)
	if info.Year != "" {
		dirName = fmt.Sprintf("%s - %s", info.Year, dirName)

	}
	if viper.GetBool("useCategory") {
		if viper.GetBool("mergeEPAndSingle") {
			if info.Category == "EP" || info.Category == "Single" {
				return path.Join(baseDir, artist, dirName, "Singles And EPs")
			}
			return path.Join(baseDir, artist, dirName, info.Category)
		}

		return path.Join(baseDir, artist, dirName, info.Category)
	}

	return path.Join(baseDir, artist, dirName)
}

// PrepareDownload 准备目录， 将下载数据发送到Channels
func (d *Downloader) PrepareDownload(info model.AlbumInfo, baseDir string) {
	logrus.Printf("开始入队列：%s", info.Name)
	d.songChan <- downloadInfo{
		object: DownloadImage{
			DownloadUrl: info.Image,
			FileType:    "jpg",
			Name:        "cover",
		},
		downloadDir: BaseAlbumDownloadDir(baseDir, info),
	}

	for _, song := range info.MusicList {
		if song.DownloadUrl == "" {
			continue
		}

		os.MkdirAll(BaseAlbumDownloadDir(baseDir, info), fs.ModePerm)

		d.songChan <- downloadInfo{
			object: DownloadMusic{
				DownloadUrl: song.DownloadUrl,
				FileType:    "mp3",
				Name:        song.Name,
				Artist:      song.Artist,
				index:       song.Postion,
				Category:    info.Category,
			},
			downloadDir: BaseAlbumDownloadDir(baseDir, info),
		}
	}
	logrus.Printf("入队列完成：%s", info.Name)
}

func (d *Downloader) Run(group *sync.WaitGroup) {
	logrus.Printf("线程启动=>")
	context := NewContext(d.bars, d.barWaitGroup)
	defer group.Done()
	for {
		select {
		case x, ok := <-d.songChan:
			if !ok {
				logrus.Printf("线程池关闭")
				return
			}
			d.process(x, context)

		case _, ok := <-d.closeChan:
			if !ok {
				return
			}
			return
		case sig := <-d.sigChan:
			fmt.Println("接受到来自系统的信号：", sig.String())
			d.Close()
		}
	}
}
func (d *Downloader) process(x downloadInfo, contexts ...*DownloadContext){
	err := d.crawler.Download(x.object, x.downloadDir, contexts...)
	if err != nil {
		logrus.Printf("下载文件:%s 错误：%s", x.object.getFileName(), err.Error())
	}
}

