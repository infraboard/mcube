package request

import (
	"net/http"
	"strconv"
)

const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize = 20
	// DefaultPageNumber 默认页号
	DefaultPageNumber = 1
)

// NewPageRequestFromHTTP 从HTTP请求中加载分页请求
func NewPageRequestFromHTTP(req *http.Request) *PageRequest {
	qs := req.URL.Query()

	ps := qs.Get("page_size")
	pn := qs.Get("page_number")

	psUint64, _ := strconv.ParseUint(ps, 10, 64)
	pnUint64, _ := strconv.ParseUint(pn, 10, 64)

	if psUint64 == 0 {
		psUint64 = DefaultPageSize
	}
	if pnUint64 == 0 {
		pnUint64 = DefaultPageNumber
	}

	return &PageRequest{
		PageSize:   uint(psUint64),
		PageNumber: uint(pnUint64),
	}
}

// NewPageRequest 实例化
func NewPageRequest(ps uint, pn uint) *PageRequest {
	return &PageRequest{
		PageSize:   ps,
		PageNumber: pn,
	}
}

// PageRequest 分页请求 request
type PageRequest struct {
	PageSize   uint `json:"page_size,omitempty" validate:"gte=1,lte=200"`
	PageNumber uint `json:"page_number,omitempty" validate:"gte=1"`
}

// Offset skip
func (p *PageRequest) Offset() int64 {
	return int64(p.PageSize * (p.PageNumber - 1))
}
