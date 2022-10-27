package driver

import (
	"BaiduP-PST/config"
	"BaiduP-PST/pkg/base"
	"encoding/json"
	"fmt"
	"strings"
)

type quota struct {
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used"`
	Expire bool   `json:"expire"`
}

func RefreshToken() {
	var requestBody string
	tokenBody := *base.TokenBody
	setting := *config.Setting
	url := fmt.Sprintf("https://openapi.baidu.com/oauth/2.0/token?grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", tokenBody.RefreshToken, setting.AppKey, setting.SecretKey)

	bytes, err := SendHttpRequest("GET", url, strings.NewReader(requestBody), nil)

	err = json.Unmarshal(bytes, &tokenBody)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("最终结果：%s,%s,%d", tokenBody.AccessToken, tokenBody.RefreshToken, tokenBody.ExpiresIn)
}

func CheckToken() {
	var requestBody string
	tokenBody := *base.TokenBody
	url := fmt.Sprintf("https://pan.baidu.com/api/quota?access_token=%s&checkfree=1&checkexpire=1", tokenBody.AccessToken)
	fmt.Println(url)

	m := map[string]string{"User-Agent": "pan.baidu.com"}
	bytes, err := SendHttpRequest("GET", url, strings.NewReader(requestBody), m)

	if err != nil {
		fmt.Println(err)
	}

	q := new(quota)
	err = json.Unmarshal(bytes, q)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("校验成功，查询网盘容量结果:总共 %d GB,已用 %d GB", accGB(q.Total), accGB(q.Used))
}

func accGB(n uint64) (f uint64) {
	f = n / 1024 / 1024 / 1024
	return
}
