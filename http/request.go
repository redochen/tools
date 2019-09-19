package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

//HTTP请求类
type Request struct {
	Url                string            // 要请求的ULR地址
	Method             string            // 请求类型：GET，POST，默认：GET
	PostData           map[string]string // 请求要发送的数据
	PostString         string            // 请求要发送的字符串
	SetOpaque          bool              // 是否设置URL不转码
	ContentEncoding    string            // 提交数据：PostData_Text 的编码，默认值：UTF8
	Accept             string            // Accept: application/json
	UserName           string            // 用户账号
	Password           string            // 用户密码
	AcceptEncoding     string            // 声明浏览器支持的编码类型:gzip,deflate,identity 中的一种或多种，逗号分开。默认：gzip,deflate
	AcceptLanguage     string            // 支持的语言
	ContentType        string            // 内容类型，默认值： "application/x-www-form-urlencoded; charset=UTF-8";
	TimeoutSeconds     int               // 请求超时时间（秒为单位），默认：30秒
	WebProxy           string            // 代理服务器IP地址：222.222.222.222:80
	BypassProxyOnLocal bool              // 如果对本地地址不使用代理服务器，则为 true；否则为 false。默认值为 false。
	UserAgent          string            // 浏览器信息，默认：	Mozilla/5.0 (Windows NT 6.1; WOW64; rv:20.0) Gecko/20100101 Firefox/20.0
	Referer            string            // 引用页 默认值string.Empty。
	Origin             string            // 来源
	Headers            map[string]string // 头信息
	ErrorInfo          string            // 错误信息
	AllowAutoRedirect  bool              // 是否根据301跳转
	Host               string            // 服务器主机
	Cookies            []*http.Cookie    // Cookies
	CookieJar          *cookiejar.Jar    // CookiesJar
}

//发送GET请求
func (r *Request) Get() (string, error) {
	return r.Do()
}

//发送GET请求
func (r *Request) GetEx() (string, []*http.Cookie, error) {
	return r.DoEx()
}

//发送POST请求
func (r *Request) Post() (string, error) {
	return r.Do()
}

//发送POST请求
func (r *Request) PostEx() (string, []*http.Cookie, error) {
	return r.DoEx()
}

//发送请求
func (r *Request) Do() (string, error) {
	rsp, _, err := r.DoEx()
	return rsp, err
}

//发送请求
func (r *Request) DoEx() (string, []*http.Cookie, error) {
	if r.TimeoutSeconds > 0 {
		DialTimeoutSeconds = time.Duration(r.TimeoutSeconds)
	}

	client := &http.Client{
		CheckRedirect: CheckRedirect,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext:     DialContext,
		},
	}

	if r.CookieJar != nil {
		client.Jar = r.CookieJar
	}

	var reader io.Reader

	//转换POST数据
	if r.PostData != nil && len(r.PostData) > 0 {
		qs := GetQueryString(r.PostData)
		reader = strings.NewReader(qs)
	} else if r.PostString != "" {
		reader = strings.NewReader(r.PostString)
	} else {
		reader = nil
	}

	req, err := NewHttpRequest(r, reader)
	if err != nil {
		return "", nil, err
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode != 200 {
		return fmt.Sprintf("StatusCode=%d", resp.StatusCode), nil, err
	} else {
		body, err := GetBody(resp)
		//cookies, _ := GetCookies(resp)
		if r.CookieJar != nil {
			cookies := r.CookieJar.Cookies(req.URL)
			return body, cookies, err
		} else {
			return body, nil, err
		}
	}
}
