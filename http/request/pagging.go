package request

// PageRequet 分页请求
type PageRequet struct {
	PageSize   uint `json:"page_size,omitempty" validate:"gte=1,lte=200"`
	PageNumber uint `json:"page_number,omitempty" validate:"gte=1"`
}
