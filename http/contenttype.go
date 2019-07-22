package http

type ContentType int

const (
	TextPlain ContentType = iota
	TextHtml
	TextXml
	ApplicationJson
	ApplicationXml
	ApplicationXWwwFormUrlEncoded
	MultipartFormData
	MultipartMixed
	ImageGif
	ImageJpeg

	UnknownType //请在此之前添加新的定义
)

//获取ContentType的字符串值
func (ct ContentType) String() string {
	switch ct {
	case TextPlain:
		return "text/plain"
	case TextHtml:
		return "text/html"
	case TextXml:
		return "text/xml"
	case ApplicationJson:
		return "application/json"
	case ApplicationXml:
		return "application/xml"
	case ApplicationXWwwFormUrlEncoded:
		return "application/x-www-form-urlencoded"
	case MultipartFormData:
		return "multipart/form-data"
	case MultipartMixed:
		return "multipart/mixed"
	case ImageGif:
		return "image/gif"
	case ImageJpeg:
		return "image/jpeg"
	}
	return "UnknownType"
}

//获取ContentType的整数值
func (ct ContentType) Value() int {
	if ct >= TextPlain && ct < UnknownType {
		return int(ct)
	} else {
		return int(UnknownType)
	}
}
