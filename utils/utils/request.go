package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	client *http.Client
	once   sync.Once
)

func DoRequest(method string, url string, data string, headers map[string]string) (res string, err error) {
	once.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        10,               // 最大空闲连接数
				MaxIdleConnsPerHost: 10,               // 每个主机的最大空闲连接数
				IdleConnTimeout:     30 * time.Second, // 空闲连接超时时间
			},
		}
	})

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return string(body), nil
}
