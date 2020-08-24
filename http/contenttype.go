package http

//ContentType ContentType枚举定义
type ContentType int

const (
	//TextPlain 普通文本
	TextPlain ContentType = iota

	//TextHTML Html文本
	TextHTML

	//TextXML Xml文本
	TextXML

	//ApplicationJSON Json内容
	ApplicationJSON

	// ApplicationXML Xml内容
	ApplicationXML

	//ApplicationXWwwFormURLEncoded URL编码后的表单内容
	ApplicationXWwwFormURLEncoded

	//MultipartFormData 多分部表单数据
	MultipartFormData

	//MultipartMixed 多分部混合数据
	MultipartMixed

	//ImageGif GIF图片
	ImageGif

	//ImageJpeg JPEG图片
	ImageJpeg

	//UnknownType 未知类型
	UnknownType //请在此之前添加新的定义
)

//获取ContentType的字符串值
func (ct ContentType) String() string {
	switch ct {
	case TextPlain:
		return "text/plain"
	case TextHTML:
		return "text/html"
	case TextXML:
		return "text/xml"
	case ApplicationJSON:
		return "application/json"
	case ApplicationXML:
		return "application/xml"
	case ApplicationXWwwFormURLEncoded:
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

//Value 获取ContentType的整数值
func (ct ContentType) Value() int {
	if ct >= TextPlain && ct < UnknownType {
		return int(ct)
	}

	return int(UnknownType)
}
