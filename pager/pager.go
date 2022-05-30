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
	Scan(context.Context, Set) error
	// 设置页面打小, 默认20, 一页数据20条
	SetPageSize(ps int64)
	// 设置读取速率, 默认1, 每秒发起一次请求
	SetRate(r float64)
}

// 可以往里面添加元素
type Set interface {
	// 往Set添加元素
	Add(...any)
	// 当前Set有多少个元素
	Length() int64
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

func (p *BasePager) PageSize() int64 {
	return p.pageSize
}

func (p *BasePager) PageNumber() int64 {
	return p.pageNumber
}

func (p *BasePager) Next() bool {
	p.tb.Wait(1)
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

// 通过判断当前set是否小于PageSize, 从而判断是否满页
func (p *BasePager) CheckHasNext(set Set) {
	if int64(set.Length()) < p.pageSize {
		p.hasNext = false
	} else {
		p.pageNumber++
	}
}
