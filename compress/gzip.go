package compress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
)

//BestZip 压缩字符串（最好压缩比例）
func BestZip(s string) ([]byte, error) {
	return ZipData([]byte(s), gzip.BestCompression)
}

//FastZip 压缩字符串（最快压缩速度）
func FastZip(s string) ([]byte, error) {
	return ZipData([]byte(s), gzip.BestSpeed)
}

//ZipData 压缩数据
func ZipData(data []byte, level int) ([]byte, error) {
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

//Unzip 解压缩字符串
func Unzip(data []byte) (string, error) {
	b, err := UnzipData(data)
	if b != nil && len(b) > 0 {
		return string(b), nil
	}

	return "", err
}

//UnzipData 解压缩数据
func UnzipData(data []byte) ([]byte, error) {
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
