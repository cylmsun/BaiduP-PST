package model

import (
	"io/fs"
)

type DicInfo struct {
	Path    string `json:"path"`
	Name    string `json:"server_filename"`
	ModTime int64  `json:"server_mtime"`
	IsDir   int8   `json:"isdir"`
	RWMode  fs.FileMode
}
