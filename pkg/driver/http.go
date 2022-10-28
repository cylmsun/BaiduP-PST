package driver

import (
	"io"
	"net/http"
)

func SendHttpRequest(method string, url string, body io.Reader, headerMap map[string]string) (bytes []byte, err error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if headerMap != nil {
		for k, v := range headerMap {
			request.Header.Set(k, v)
		}
	}

	client := http.Client{}
	do, err := client.Do(request)
	// 随手close是个好习惯
	defer do.Body.Close()
	if err != nil {
		return
	}

	bytes, err = io.ReadAll(do.Body)
	if err != nil {
		return
	}
	return
}
