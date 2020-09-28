package response

// Data to be used by controllers.
type Data struct {
	Code      *int        `json:"code"`                // 自定义返回码  0:表示正常
	Type      string      `json:"type,omitempty"`      // 数据类型, 可以缺省
	Namespace string      `json:"namespace,omitempty"` // 异常的范围
	Reason    string      `json:"reason,omitempty"`    // 异常原因
	Message   string      `json:"message,omitempty"`   // 关于这次响应的说明信息
	Data      interface{} `json:"data,omitempty"`      // 返回的具体数据
	Meta      interface{} `json:"meta,omitempty"`      // 数据meta
}

// NewData new实例
func NewData(data interface{}) *Data {
	code := -1
	return &Data{
		Code: &code,
		Data: data,
	}
}
