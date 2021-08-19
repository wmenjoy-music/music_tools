package main

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"wmenjoy/music/pkg/app"
)

func main() {
	logrus.SetOutput(colorable.NewColorableStdout())
	if err := app.MainErr(); err != nil {
		logrus.Fatal(err)
	}
}
