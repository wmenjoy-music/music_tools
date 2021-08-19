package service

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
	"wmenjoy/music/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
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
	proxy       string
	Retry       int
	Options     Options
	ProcessBars *mpb.Progress
	count       int
	bars        []*mpb.Bar
}

// ParsePage 使用get方法获取页面
func (c *Crawler) ParsePage(url string, objectConsumer func(Body io.Reader) (interface{}, error)) (interface{}, error) {
	rand := time.Duration(rand.Intn(2))
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

func (c *Crawler) Download(obj IDownloadObject, downloadDir string, contexts ...*DownloadContext) error {
	for count := 0; count <= c.Retry; count++ {
		err := c.__download(obj, downloadDir, contexts...)
		if err == nil {
			return nil
		}
		logrus.Printf("下载文件:%s 错误：%s", obj.getFileName(), err.Error())
	}

	return nil
}

type barWriter struct {
	io.Writer
	bar   *mpb.Bar
	start time.Time
	count int
}

// Write implement io.Writer
func (p *barWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	p.count += n

	p.bar.IncrBy(n)
	p.bar.DecoratorEwmaUpdate(time.Since(p.start))
	return
}

func (p *barWriter) Close() (err error) {
	return
}

type DownloadContext struct {
	Index int
	Bars  []*mpb.Bar
}

var count = -1

var Bars []*mpb.Bar

func NewContext(bars []*mpb.Bar, barWaitGroup []*sync.WaitGroup) *DownloadContext {
	count++
	return &DownloadContext{
		Index: count,
		Bars:  bars,
	}
}

func (c *Crawler) __download(obj IDownloadObject, downloadDir string, contexts ...*DownloadContext) error {
	if obj == nil {
		return errors.New("不合法的下载对象")
	}

	fileName := utils.ValidateFileName(obj.getFileName())
	logrus.Printf("开始下载文件：%s", path.Join(downloadDir, fileName))
	if val, _ := utils.PathExists(path.Join(downloadDir, fileName)); val {
		logrus.Printf("文件已经下载：%s", path.Join(downloadDir, fileName))
		return nil
	}
	rand := time.Duration(rand.Intn(500))
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

		if len(contexts) > 0 {

		}
		logrus.Printf("下载完成文件：%s", fileName)

	}(f)

	if err != nil {
		return err
	}
	out = f

	if c.Options.ShowProgress {
		job := fileName[strings.Index(fileName, "-")+2:]
		task := "downloading"
		var bar *mpb.Bar
		if len(contexts) == 0 || contexts[0].Bars[contexts[0].Index] == nil {
			if len(contexts) > 0 {
				task = task + "-" + strconv.Itoa(contexts[0].Index)
			}

			bar = c.ProcessBars.AddBar(length,
				//	mpb.BarFillerClearOnComplete(),
				mpb.BarRemoveOnComplete(),
				mpb.PrependDecorators(
					// simple name decorator
					decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
					decor.Name(job, decor.WCSyncSpaceR),
					// decor.DSyncWidth bit enables column width synchronization
					decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 5}),
					// replace ETA decorator with "done" message, OnComplete event
					decor.OnComplete(
						// ETA decorator with ewma age of 60
						decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncWidth), "done"),
				),
			)
			if len(contexts) > 0 {
				contexts[0].Bars[contexts[0].Index] = bar
			}
		} else {
			if len(contexts) > 0 {
				logrus.Printf("----->%d", contexts[0].Index)
				task = task + "-" + strconv.Itoa(contexts[0].Index)
			}
			bar = c.ProcessBars.AddBar(length,
				//mpb.BarQueueAfter(contexts[0].Bars[contexts[0].Index]),
				mpb.BarRemoveOnComplete(),
				//	mpb.BarFillerClearOnComplete(),
				mpb.PrependDecorators(
					// simple name decorator
					decor.Name(task, decor.WC{W: len(task) + 1, C: decor.DidentRight}),
					decor.Name(job, decor.WCSyncSpaceR),
					// decor.DSyncWidth bit enables column width synchronization
					decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 5}),
					// replace ETA decorator with "done" message, OnComplete event
					decor.OnComplete(
						// ETA decorator with ewma age of 60
						decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncWidth), "done",
					),
				),
			)
			contexts[0].Bars[contexts[0].Index] = bar
		}
		/*

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
				}))*/
		out = io.MultiWriter(out, &barWriter{
			bar:   bar,
			start: time.Now(),
		})
	}

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		logrus.Printf("下载完成文件：%s 失败:%s", fileName, err.Error())
		return err
	}

	return err
}
