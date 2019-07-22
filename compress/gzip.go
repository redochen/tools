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
	return h.ZipData([]byte(s), gzip.BestCompression)
}

//压缩字符串（最快压缩速度）
func (h *GZipHelper) FastZip(s string) ([]byte, error) {
	return h.ZipData([]byte(s), gzip.BestSpeed)
}

//压缩数据
func (h *GZipHelper) ZipData(data []byte, level int) ([]byte, error) {
	if nil == data || len(data) <= 0 {
		return nil, errors.New("input buffer is nil or empty")
	}

	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, level)
	if err != nil {
		return nil, err
	}

	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

//解压缩字符串
func (h *GZipHelper) Unzip(data []byte) (string, error) {
	b, err := h.UnzipData(data)
	if b != nil && len(b) > 0 {
		return string(b), nil
	} else {
		return "", err
	}
}

//解压缩数据
func (h *GZipHelper) UnzipData(data []byte) ([]byte, error) {
	if nil == data || len(data) <= 0 {
		return nil, errors.New("input buffer is nil or empty")
	}

	b := bytes.NewBuffer(data)

	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	return ioutil.ReadAll(r)
}
