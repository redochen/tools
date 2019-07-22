package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	CcHttp = NewHttpHelper()
)

//HTTP帮助类
type HttpHelper struct {
}

//获取一个新的HttpHelper实例
func NewHttpHelper() *HttpHelper {
	return &HttpHelper{}
}

//GET方式请求URL
func (h *HttpHelper) Get(url string) (string, error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), err
	} else {
		body, err := GetBody(resp)
		return body, err
	}
}

//GET方式请求URL
func (h *HttpHelper) GetEx(url, parameter string) (string, error) {
	return h.Get(fmt.Sprintf("%s%s", url, parameter))
}

//POST方式请求URL
func (h *HttpHelper) Post(url, parameter string, contentType ContentType) (string, error) {
	reader := bytes.NewBufferString(parameter)
	return h.PostEx(url, reader, contentType)
}

//POST方式请求URL
func (h *HttpHelper) PostEx(url string, reader io.Reader, contentType ContentType) (string, error) {
	resp, err := http.Post(url, contentType.String(), reader)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), err
	} else {
		body, err := GetBody(resp)
		return body, err
	}
}
