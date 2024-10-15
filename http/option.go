package http

import (
	"bytes"
	"io"
	"time"
)

// HttpOption HTTP选项类
type HttpOption struct {
	Parameter        string        //URL参数或表单参数
	ContentType      ContentType   //参数类型
	Timeout          time.Duration //超时
	IgnoreStatusCode bool          //是否忽略HTTP状态
}

// DefaultHttpOption 获取新的HttpOption实例
func DefaultHttpOption() *HttpOption {
	return NewHttpOptionEx("", TextPlain, DefaultTimeout, false)
}

// NewHttpOption 获取新的HttpOption实例
func NewHttpOption(parameter string, contentType ContentType) *HttpOption {
	return NewHttpOptionEx(parameter, contentType, DefaultTimeout, false)
}

// NewHttpOptionEx 获取新的HttpOption实例
func NewHttpOptionEx(parameter string, contentType ContentType, timeout time.Duration, ignoreStatusCode bool) *HttpOption {
	return &HttpOption{
		Parameter:        parameter,
		ContentType:      contentType,
		Timeout:          timeout,
		IgnoreStatusCode: ignoreStatusCode,
	}
}

// EnsureHttpOption 确保HttpOption有效
func EnsureHttpOption(option *HttpOption) *HttpOption {
	if option == nil {
		option = DefaultHttpOption()
	}

	if option.Timeout == 0 {
		option.Timeout = DefaultTimeout
	}

	return option
}

// GetParameterReader 获取参数读取器实例
func (option *HttpOption) GetParameterReader() io.Reader {
	if option.Parameter == "" {
		return nil
	}
	return bytes.NewBufferString(option.Parameter)
}
