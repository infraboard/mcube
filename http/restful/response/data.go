package response

import "github.com/infraboard/mcube/http/response"

// NewData new实例
func NewData(data interface{}) *response.Data {
	return response.NewData(data)
}
