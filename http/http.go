package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

//Get GET方式请求URL
func Get(url string) (string, error) {
	resp, err := http.Get(url)
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

//GetEx GET方式请求URL
func GetEx(url, parameter string) (string, error) {
	return Get(fmt.Sprintf("%s%s", url, parameter))
}

//Post POST方式请求URL
func Post(url, parameter string, contentType ContentType) (string, error) {
	reader := bytes.NewBufferString(parameter)
	return PostEx(url, reader, contentType)
}

//PostEx POST方式请求URL
func PostEx(url string, reader io.Reader, contentType ContentType) (string, error) {
	resp, err := http.Post(url, contentType.String(), reader)
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
