package util

import (
	"errors"
	"fmt"
	"os"
)

// CheckDir 检查路径目录,不实现创建目录的功能
func CheckDir(path *string) (b bool, err error) {
	b = false
	stat, err := os.Stat(*path)
	if os.IsNotExist(err) {
		err = errors.New(fmt.Sprintf("no such dictionary:%s", *path))
		return
	}
	if !stat.IsDir() {
		err = errors.New(fmt.Sprintf("%s is not a dictionary", *path))
		return
	}
	b = true
	return
}
