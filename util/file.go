package util

import (
	"BaiduP-PST/pkg/base"
	"BaiduP-PST/pkg/driver"
	"BaiduP-PST/pkg/model"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type fileSplit struct {
	num int
	f   *os.File
	md5 string
}

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
	//1. slice file
	fileSpits := splitFile(path)
	// 2.precreate
	preCreateId := preCreate(path, fileSpits)
	// 3.upload
	splitsUpload(fileSpits, preCreateId)
	// 4.join
}

// todo order splits
func splitFile(path string) (f []fileSplit) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	if err != nil {
		fmt.Printf("readDir error:%s", err.Error())
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("fileInfo error:%s", err.Error())
		return
	}
	splitNums := (int(fileInfo.Size()) + bufferSize - 1) / bufferSize

	chanF := make(chan fileSplit, splitNums)
	defer close(chanF)
	for i := 0; i < splitNums; i++ {
		go createTmpFile(file.Name(), fileInfo.Name(), "", i, chanF)
	}

	tempFiles := make([]fileSplit, splitNums, splitNums)
	// 期望：协程全部结束
	for i := 0; i < splitNums; i++ {
		buffer := make([]byte, bufferSize)
		tempFiles[i] = <-chanF
		at, err := file.ReadAt(buffer, int64(bufferSize*tempFiles[i].num))
		if err != nil {
			fmt.Printf("ReadAt error %d:%s", at, err.Error())
			return
		}
		tempFiles[i].f.Write(buffer)

		md5String := md5.New()
		toString := hex.EncodeToString(md5String.Sum(buffer))
		tempFiles[i].md5 = toString
	}

	return tempFiles

}

func createTmpFile(path string, name string, suffix string, num int, c chan<- fileSplit) {
	// fileName: ***/name_splitNumber.suffix
	tempName := fmt.Sprintf("%s/%s_%d.%s", path, name, num, suffix)
	f, err := os.Create(tempName)
	if err != nil {
		fmt.Printf("createTmpFile error:%s", err.Error())
	}
	spit := fileSplit{f: f, num: num}
	c <- spit
}

func preCreate(path string, f []fileSplit) (id string) {
	var requestBody map[string]any
	tokenBody := *base.TokenBody
	u := fmt.Sprintf("https://pan.baidu.com/rest/2.0/xpan/file?method=precreate&access_token=%s", tokenBody.AccessToken)
	m := map[string]string{"User-Agent": "pan.baidu.com"}

	blockList := make([]string, len(f), len(f))
	for i := 0; i < len(f); i++ {
		blockList[i] = f[i].md5
	}

	requestBody["path"] = path
	requestBody["size"] = 0
	requestBody["isdir"] = 0
	requestBody["autoinit"] = 1
	requestBody["block_list"] = blockList

	j, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("170 error:%s", err.Error())
	}
	responseBody, err := driver.SendHttpRequest("POST", u, strings.NewReader(string(j)), m)
	if err != nil {
		fmt.Printf("precreate error:%s \n", err.Error())
		return
	}

	var res model.PreCreateResp
	err = json.Unmarshal(responseBody, &res)
	if err != nil {
		fmt.Printf("SendHttpRequest json.Unmarshal error:%s \n", err.Error())
		return
	}
	if res.ErrNo != 0 {
		fmt.Printf("precreate error:%s \n", err.Error())
	}
	id = res.UploadId
	return id
}

func splitsUpload(f []fileSplit, id string) {
	chanF := make(chan string, len(f))
	for i := 0; i < len(f); i++ {
		go upload(f[i], id, chanF)
	}
}

func upload(f fileSplit, id string, chanF chan<- string) {
	stat, _ := f.f.Stat()
	size := stat.Size()
	bytes := make([]byte, size)
	f.f.ReadAt(bytes, 0)

	tokenBody := *base.TokenBody
	u := fmt.Sprintf("http://d.pcs.baidu.com/rest/2.0/pcs/superfile2?method=upload&access_token=%s", tokenBody.AccessToken)
	m := map[string]string{}

	//_, err := driver.SendHttpRequest("POST", u, strings.NewReader(string(bytes)), m)
	//if err != nil {
	//	fmt.Printf("upload split file error:%s", err.Error())
	//}

}

func bool2int(b bool) (i int8) {
	if b {
		return 1
	} else {
		return 0
	}
}
