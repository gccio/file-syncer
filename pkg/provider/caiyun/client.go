package caiyun

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func NewRequest(path, method string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, basePath+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	key := getRandomStr(16)
	sign := getSign(ts, key, string(body))

	req.Header.Add("Cookie", cookie)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("x-huawei-channelSrc", "10000034")
	req.Header.Add("x-inner-ntwk", "2")
	req.Header.Add("mcloud-channel", "1000101")
	req.Header.Add("mcloud-client", "10701")
	req.Header.Add("mcloud-sign", fmt.Sprintf("%s,%s,%s", ts, key, sign))
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("caller", "web")
	req.Header.Add("CMS-DEVICE", "default")
	req.Header.Add("x-DeviceInfo", "||9|85.0.4183.83|chrome|85.0.4183.83|||windows 10||zh-CN|||")
	req.Header.Add("x-SvcType", "1")
	req.Header.Add("referer", "https://yun.139.com/w/")

	return req, nil
}
