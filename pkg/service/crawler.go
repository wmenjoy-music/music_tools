package service

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
	"wmenjoy/music/pkg/utils"

	ansi "github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

type Options struct {
	ShowProgress  bool
	Debug         bool
	DoNotDownload bool
	FileName      string
}

type IDownloadObject interface {
	//download
	getDownloadUrl() string
	getFileName() string
}

type DownloadImage struct {
	Name        string
	FileType    string
	DownloadUrl string
}

var _ IDownloadObject = (*DownloadImage)(nil)

func (obj DownloadImage) getFileName() string {
	return fmt.Sprintf("%s.%s", obj.Name, obj.FileType)
}

func (obj DownloadImage) getDownloadUrl() string {
	return obj.DownloadUrl
}

type DownloadMusic struct {
	IDownloadObject
	index       string
	Name        string
	Artist      string
	Category    string
	FileType    string
	DownloadUrl string
}

var _ IDownloadObject = (*DownloadMusic)(nil)

func (obj DownloadMusic) getFileName() string {
	return fmt.Sprintf("%s. %s - %s.%s", obj.index, obj.Artist, obj.Name, obj.FileType)
}

func (obj DownloadMusic) getDownloadUrl() string {
	return obj.DownloadUrl
}

type Crawler struct {
	proxy   string
	Retry   int
	Options Options
}

// ParsePage 使用get方法获取页面
func (c Crawler) ParsePage(url string, objectConsumer func(Body io.Reader) (interface{}, error)) (interface{}, error) {
	rand := time.Duration(rand.Intn(1000))
	time.Sleep(rand * time.Millisecond)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return objectConsumer(res.Body)
}

func (c Crawler) Download(obj IDownloadObject, downloadDir string) error {
	if obj == nil {
		return errors.New("不合法的下载对象")
	}

	fileName := utils.ValidateFileName(obj.getFileName())
	logrus.Printf("开始下载文件：%s", path.Join(downloadDir, fileName))
	if val, _ := utils.PathExists(path.Join(downloadDir, fileName)); val {
		logrus.Printf("文件已经下载：%s", path.Join(downloadDir, fileName))
		return nil
	}
	rand := time.Duration(rand.Intn(1000))
	time.Sleep(rand * time.Millisecond)

	resp, err := http.Get(obj.getDownloadUrl())

	if err != nil {
		return err
	}

	length := resp.ContentLength
	defer resp.Body.Close()

	var out io.Writer

	f, err := os.Create(path.Join(downloadDir, fileName+".bak"))
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Printf("close File %s Error:%s", path.Join(downloadDir, fileName, ".bak"), err.Error())
		}
		err = os.Rename(path.Join(downloadDir, fileName+".bak"), path.Join(downloadDir, fileName))

		if err != nil {
			logrus.Printf("下载完成文件：%s 失败:%s", fileName, err.Error())
		}

		logrus.Printf("下载完成文件：%s", fileName)
	}(f)

	if err != nil {
		return err
	}
	out = f

	if c.Options.ShowProgress {
		bar := progressbar.NewOptions64(length,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription(fileName),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))
		out = io.MultiWriter(out, bar)
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		logrus.Printf("下载完成文件：%s 失败:%s", fileName, err.Error())
		return err
	}

	return err
}
