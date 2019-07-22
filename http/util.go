package http

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	. "github.com/redochen/tools/crypto"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	DialTimeoutSeconds time.Duration = 30 //超时秒数
)

//HTML解码
func HtmlDecode(html string) string {

	//解码<
	lt, _ := regexp.Compile(`&lt;`)
	html = lt.ReplaceAllString(html, "<")

	//解码>
	gt, _ := regexp.Compile(`&gt;`)
	html = gt.ReplaceAllString(html, ">")

	//解码"
	slat, _ := regexp.Compile(`&#034;`)
	html = slat.ReplaceAllString(html, "\"")

	return html
}

//URL编码
func UrlEncode(path string) string {
	u := &url.URL{Path: path}
	return u.String()
}

//创建新的HTTP请求（重写http.NewRequest方法，因为该方法会对URL进行转码）
func NewHttpRequest(r *Request, body io.Reader) (*http.Request, error) {
	if nil == r {
		return nil, errors.New("request is nil")
	}

	method := r.Method
	if len(method) <= 0 {
		method = "GET"
	}

	//httpReq, err := http.NewRequest(method, r.Url, body)
	httpReq, err := NewRawRequest(method, r.Url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Connection", "keep-alive")

	if len(r.Method) > 0 {
		httpReq.Header.Add("Method", method)
	} else {
		httpReq.Header.Add("Method", "GET")
	}

	if len(r.Accept) > 0 {
		httpReq.Header.Add("Accept", r.Accept)
	}

	if len(r.AcceptEncoding) > 0 {
		httpReq.Header.Add("Accept-Encoding", r.AcceptEncoding)
	} else {
		httpReq.Header.Add("Accept-Encoding", "gzip,deflate")
	}

	if len(r.UserAgent) > 0 {
		httpReq.Header.Add("User-Agent", r.UserAgent)
	}

	if len(r.ContentType) > 0 {
		httpReq.Header.Set("Content-Type", r.ContentType)
	} else {
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	}

	auth, err := HttpsBasicAuthorization(r.UserName, r.Password)
	if err != nil {
		return nil, err
	} else if len(auth) > 0 {
		httpReq.Header.Add("Authorization", auth)
	}

	if len(r.Referer) > 0 {
		httpReq.Header.Add("Referer", r.Referer)
	}

	if len(r.Origin) > 0 {
		httpReq.Header.Set("Origin", r.Origin)
	}

	if len(r.Host) > 0 {
		httpReq.Header.Set("Host", r.Host)
	}

	if r.Headers != nil {
		for k, v := range r.Headers {
			httpReq.Header.Set(k, v)
		}
	}

	return httpReq, nil
}

//创建新的HTTP请求（重写http.NewRequest方法，因为该方法会对URL进行转码）
func NewRawRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	httpUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	//设置不转码路径
	httpUrl.Opaque = urlStr

	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}

	httpReq := &http.Request{
		Method:     method,
		URL:        httpUrl,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Body:       rc,
		Host:       httpUrl.Host,
	}

	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			httpReq.ContentLength = int64(v.Len())
		case *bytes.Reader:
			httpReq.ContentLength = int64(v.Len())
		case *strings.Reader:
			httpReq.ContentLength = int64(v.Len())
		}
	}

	return httpReq, nil
}

//获取原始请求参数
func GetRawParameter(r *http.Request) (string, error) {
	if nil == r {
		return "", errors.New("request is nil")
	}

	// if we would use this, then the POST and GET requests would be merged
	// to the req.Form variable
	// req.ParseForm()
	// v := req.Form

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

//获取HTTP响应体
func GetBody(r *http.Response) (string, error) {
	if nil == r {
		return "", errors.New("response is nil")
	}

	var body string
	encoding := r.Header.Get("Content-Encoding")
	switch encoding {
	case "gzip":
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()

		for {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)

			if err != nil && err != io.EOF {
				return "", err
			}

			if n == 0 {
				break
			}

			body += string(buf[:n])
		}

	default:
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", err
		}
		body = string(bodyBytes)
	}

	return body, nil
}

//获取Cookies
func GetCookies(r *http.Response) ([]*http.Cookie, error) {
	if nil == r {
		return nil, errors.New("response is nil")
	} else {
		return r.Cookies(), nil
	}
}

//打印Cookies
func PrintCookies(cookies []*http.Cookie, page string) {
	if cookies == nil {
		println(fmt.Sprintf("[%s] cookies is null", page))
	} else {
		var length = len(cookies)
		println(fmt.Sprintf("[%s] cookies count: %d", page, length))

		for i := 0; i < length; i++ {
			var c = cookies[i]

			println(fmt.Sprintf("------ Cookie [%d]------", i))
			println(fmt.Sprintf("Name\t=%s", c.Name))
			println(fmt.Sprintf("Value\t=%s", c.Value))
			println(fmt.Sprintf("Path\t=%s", c.Path))
			println(fmt.Sprintf("Domain\t=%s", c.Domain))
			println(fmt.Sprintf("Expires\t=%s", c.Expires))
			println(fmt.Sprintf("RawExpires=%s", c.RawExpires))
			println(fmt.Sprintf("MaxAge\t=%d", c.MaxAge))
			println(fmt.Sprintf("Secure\t=%t", c.Secure))
			println(fmt.Sprintf("HttpOnly=%t", c.HttpOnly))
			println(fmt.Sprintf("Raw\t=%s", c.Raw))
			println(fmt.Sprintf("Unparsed=%s", c.Unparsed))
		}
	}
}

//检测重定向
func CheckRedirect(r *http.Request, via []*http.Request) error {
	println(fmt.Sprintf("After %d redirects, the last url is %s.",
		len(via), r.URL.RequestURI()))
	return nil
}

//处理连接超时
func DialTimeout(network, addr string) (net.Conn, error) {
	deadline := time.Now().Add(DialTimeoutSeconds * time.Second)
	c, err := net.DialTimeout(network, addr, time.Duration(DialTimeoutSeconds*time.Second))
	if err != nil {
		return nil, err
	}
	c.SetDeadline(deadline)
	return c, nil
}

//创建HTTPS认证字符串
func HttpsBasicAuthorization(username, password string) (string, error) {
	if len(username) <= 0 {
		return "", nil
	}

	auth, err := CcBase64.EncodeString(username + ":" + password)
	if err != nil {
		return "", err
	}

	return "Basic " + auth, nil
}
