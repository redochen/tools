package compress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
)

var (
	CcGzip = NewGZipHelper()
)

//GZip帮助类
type GZipHelper struct {
}

//获取一个新的GZipHelper实例
func NewGZipHelper() *GZipHelper {
	return &GZipHelper{}
}

//压缩字符串（最好压缩比例）
func (h *GZipHelper) BestZip(s string) ([]byte, error) {
	return h.ZipEx([]byte(s), gzip.BestCompression)
}

//压缩字符串（最快压缩速度）
func (h *GZipHelper) FastestZip(s string) ([]byte, error) {
	return h.ZipEx([]byte(s), gzip.BestSpeed)
}

//压缩
func (h *GZipHelper) ZipEx(buffer []byte, level int) ([]byte, error) {
	if nil == buffer || len(buffer) <= 0 {
		return nil, errors.New("input buffer is nil or empty")
	}

	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, level)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	if _, err := gz.Write(buffer); err != nil {
		return nil, err
	}

	if err := gz.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

//解压缩
func (h *GZipHelper) Unzip(buffer []byte) ([]byte, error) {
	if nil == buffer || len(buffer) <= 0 {
		return nil, errors.New("input buffer is nil or empty")
	}

	b := bytes.NewBuffer(buffer)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	return ioutil.ReadAll(gz)
}
