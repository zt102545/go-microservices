package utils

import (
	"bytes"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/rest/httpc"
	"io"
	"net/http"
	"net/url"
)

func CallHttp(url string, method string, req interface{}) (string, error) {
	if url == "" || method == "" {
		return "", errors.New("url or method is empty")
	}
	resp, err := httpc.Do(context.Background(), method, url, req)

	if err != nil {
		return "", err
	}

	bodyStr, bodyStrErr := io.ReadAll(resp.Body)

	if bodyStrErr != nil {
		return "", bodyStrErr
	}

	return string(bodyStr), nil
}

func CallBodyHttp(url string, body string) (string, error) {
	if url == "" || body == "" {
		return "", errors.New("url or method is empty")
	}

	req, requestErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	if requestErr != nil {
		return "", requestErr
	}

	resp, respErr := httpc.DoRequest(req)

	if respErr != nil {
		return "", respErr
	}

	bodyStr, bodyStrErr := io.ReadAll(resp.Body)

	if bodyStrErr != nil {
		return "", bodyStrErr
	}

	return string(bodyStr), nil
}

func IsUrl(rawUrl string) bool {
	_, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return false
	}

	return true
}
