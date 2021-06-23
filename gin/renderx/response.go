package renderx

type Response struct {
	Code    int         `json:"code"`             // 业务编码
	Message string      `json:"message"`          // 错误描述
	Reason  interface{} `json:"reason,omitempty"` // 错误详细描述
	Data    interface{} `json:"data"`             // 成功时返回的数据
	ID      string      `json:"id,omitempty"`     // 当前请求的唯一ID，便于问题定位，忽略也可
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithId(id string) *Response {
	r.ID = id
	return r
}

func (r *Response) WithReason(reason interface{}) *Response {
	r.Reason = reason
	return r
}

func Success(data ...interface{}) *Response {
	resp := &Response{
		Code:    0,
		Message: "OK",
	}

	if len(data) > 0 {
		resp = resp.WithData(data[0])
	}

	return resp
}

func Error(message string, reason ...interface{}) *Response {
	resp := &Response{
		Code:    0,
		Message: message,
	}

	if len(reason) > 0 {
		resp = resp.WithReason(reason[0])
	}

	return resp
}
