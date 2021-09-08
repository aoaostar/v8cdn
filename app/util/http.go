package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func V8cdnPostForm(url string, data url.Values) string {

	contentType := "application/x-www-form-urlencoded"
	body := strings.NewReader(data.Encode())
	request, _ := http.NewRequest(http.MethodPost, url, body)
	request.Header.Set("Content-Type", contentType)
	resp, _ := http.DefaultClient.Do(request)
	//resp, err := http.Post(url, contentType, bytes.NewBuffer(marshal))
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// V8cdnPostJSON Post 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func V8cdnPost(url string, data interface{}) string {

	contentType := "application/json"
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
