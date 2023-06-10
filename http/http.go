package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

// Get GET方式请求URL（默认30秒超时）
func Get(url string) (string, error) {
	return GetEx(url, "", DefaultTimeout)
}

// Get GET方式请求URL（可设置超时）
func GetEx(url, parameter string, timeout time.Duration) (string, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	client := &http.Client{
		Timeout: timeout,
	}

	if parameter != "" {
		url = fmt.Sprintf("%s%s", url, parameter)
	}

	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), err
	}

	body, err := GetBody(resp)
	return body, err
}

// Post POST方式请求URL
func Post(url, parameter string, timeout time.Duration, contentType ContentType) (string, error) {
	reader := bytes.NewBufferString(parameter)
	return PostEx(url, reader, timeout, contentType)
}

// PostEx POST方式请求URL
func PostEx(url string, reader io.Reader, timeout time.Duration, contentType ContentType) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(url, contentType.String(), reader)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), err
	}

	body, err := GetBody(resp)
	return body, err
}
