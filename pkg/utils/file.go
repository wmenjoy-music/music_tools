package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PathExists 该项文件是否存在，通过读取上层目录, 来获取文件是否存在， os.Stat 读取文件，不区分大小写
func PathExists(file string) (bool, error) {
	//logrus.Printf(file)
	if filepath.Clean(file) == "/" {
		return true, nil
	}

	if s, _ := os.Stat(file); s != nil && s.IsDir() {
		return true, nil
	}

	dirs, err := os.ReadDir(filepath.Dir(file))

	if err != nil {
		return false, err
	}

	baseName := filepath.Base(file)
	for _, dir := range dirs {
		if dir.Name() == baseName {
			return true, nil
		}
	}

	return false, nil
}

func ValidateFileName(name string) string {
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\"", "")
	name = strings.ReplaceAll(name, "..", ".")
	return regexp.MustCompile("[*:?@#/\\<>|]").ReplaceAllString(name, "-")
}
