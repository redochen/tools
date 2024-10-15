package http

import (
	"fmt"
	"net/http"
	"time"
)

const (
	HttpGet        = "GET"
	HttpPost       = "POST"
	HttpPut        = "PUT"
	HttpDelete     = "DELETE"
	DefaultTimeout = time.Duration(30) * time.Second
)

// Get GET方式请求URL（默认30秒超时）
func Get(url string) (string, error) {
	return GetEx(url, nil)
}

// GetEx GET方式请求URL
func GetEx(url string, option *HttpOption) (string, error) {
	return DoEx(url, HttpGet, option)
}

// Post POST方式请求URL
func Post(url, parameter string, contentType ContentType) (string, error) {
	option := NewHttpOption(parameter, contentType)
	return PostEx(url, option)
}

// PostEx POST方式请求URL
func PostEx(url string, option *HttpOption) (string, error) {
	return DoEx(url, HttpPost, option)
}

// Put PUT方式请求URL
func Put(url, parameter string, contentType ContentType) (string, error) {
	option := NewHttpOption(parameter, contentType)
	return PutEx(url, option)
}

// PutEx PUT方式请求URL
func PutEx(url string, option *HttpOption) (string, error) {
	return DoEx(url, HttpPut, option)
}

// Delete DELETE方式请求URL
func Delete(url string) (string, error) {
	return DeleteEx(url, nil)
}

// DeleteEx DELETE方式请求URL
func DeleteEx(url string, option *HttpOption) (string, error) {
	return DoEx(url, HttpDelete, option)
}

// DoEx 请求URL
func DoEx(url, method string, option *HttpOption) (string, error) {
	opt := EnsureHttpOption(option)

	client := &http.Client{
		Timeout: opt.Timeout,
	}

	var req *http.Request
	var err error

	if method == HttpGet || method == HttpDelete {
		if opt.Parameter != "" {
			url = fmt.Sprintf("%s%s", url, opt.Parameter)
		}

		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return "", err
		}
	} else {
		req, err = http.NewRequest(method, url, opt.GetParameterReader())
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", opt.ContentType.String())
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil || resp == nil {
		return "", err
	}

	if !opt.IgnoreStatusCode && resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), err
	}

	body, err := GetBody(resp)
	return body, err
}
