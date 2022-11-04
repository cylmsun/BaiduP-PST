package util

import (
	"BaiduP-PST/pkg/base"
	"BaiduP-PST/pkg/driver"
	"BaiduP-PST/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var bufferSize = 4 * 1024 * 1024

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

func GetCloudFolderDetails(path string) (list []model.DicInfo, err error) {
	var requestBody string
	tokenBody := *base.TokenBody
	u := fmt.Sprintf("https://pan.baidu.com/rest/2.0/xpan/multimedia?method=listall&path=%s&access_token=%s&recursion=1", url.QueryEscape(path), tokenBody.AccessToken)
	m := map[string]string{"User-Agent": "pan.baidu.com"}

	responseBody, err := driver.SendHttpRequest("GET", u, strings.NewReader(requestBody), m)
	if err != nil {
		fmt.Printf("SendHttpRequest error:%s \n", err.Error())
		return
	}

	var res model.NetDicInfoResp
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		fmt.Printf("SendHttpRequest json.Unmarshal error:%s \n", err.Error())
		return
	}

	list = res.List
	return
}

func GetLocalFolderDetails(path string) (infos []model.DicInfo) {
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("readDir error:%s", err.Error())
	}
	if stat.IsDir() {
		dir, _ := os.ReadDir(path)
		for _, entry := range dir {
			d := GetLocalFolderDetails(path + "/" + entry.Name())
			infos = append(infos, d...)
		}
	} else {
		info := model.DicInfo{
			Path:    path + "/" + stat.Name(),
			Name:    stat.Name(),
			ModTime: stat.ModTime().Unix(),
			IsDir:   bool2int(stat.IsDir()),
			RWMode:  stat.Mode(),
		}
		infos = append(infos, info)
	}
	return
}

func Upload(path string) {
	// step:
	// 1. slice file
	// 2.precreate
	// 3.superfile2
	// 4.create
}

func splitFile(path string) {
	//file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	//defer file.Close()
	//if err != nil {
	//	fmt.Printf("readDir error:%s", err.Error())
	//	return
	//}
	//
	//split_list := make([]string, 0)
	//content_md5 := ""
	//slice_md5 := ""
	//
	//data := make([]byte, 1024, 1024)

}

func createTmpFile(path string) {

}

func bool2int(b bool) (i int8) {
	if b {
		return 1
	} else {
		return 0
	}
}
