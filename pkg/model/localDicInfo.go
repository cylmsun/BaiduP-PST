package model

import (
	"io/fs"
)

type DicInfo struct {
	Path    string
	Name    string
	ModTime int64
	IsDir   int8
	RWMode  fs.FileMode
}
