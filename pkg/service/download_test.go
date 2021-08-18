package service

import (
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestCrawler_DownloadObj1(t *testing.T) {

	wg := sync.WaitGroup{}

	wg.Wait()
	logrus.Printf("jieshu")
}

func test(wg *sync.WaitGroup) {

	wg.Add(2)

	for i := 0; i < 2; i++ {

		go func() {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()
	}
}
