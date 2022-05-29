package pager

import (
	"context"

	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
)

// 分页迭代器
// for p.Next() {
// 	p.Scan(ctx, dst);
// }
type Pager interface {
	// 判断是否还有下一页
	Next() bool
	// 获取当前页面的数据, 注意必须先调用Next, 从而判断是否存在下一页
	Scan(context.Context, any) error
	// 设置页面打小, 默认20, 一页数据20条
	SetPageSize(ps int64)
	// 设置读取速率, 默认1, 每秒发起一次请求
	SetRate(r float64)
}

func NewBasePager() *BasePager {
	return &BasePager{
		pageSize:   20,
		pageNumber: 1,
		hasNext:    true,
		tb:         tokenbucket.NewBucketWithRate(1, 1),
	}
}

type BasePager struct {
	pageSize   int64
	pageNumber int64
	hasNext    bool
	tb         *tokenbucket.Bucket
}

func (p *BasePager) Next() bool {
	return p.hasNext
}

func (p *BasePager) SetPageSize(ps int64) {
	p.pageSize = ps
}

func (p *BasePager) SetRate(r float64) {
	p.tb.SetRate(r)
}

func (p *BasePager) Offset() int64 {
	return int64(p.pageSize * (p.pageNumber - 1))
}

func (p *BasePager) SetHasNext(hn bool) {
	p.hasNext = hn
}

func (p *BasePager) SetPageNumber(pn int64) {
	p.pageNumber = pn
}

func (p *BasePager) IncrPageNumber(pn int64) {
	p.pageNumber++
}
