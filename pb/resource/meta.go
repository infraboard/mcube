package resource

import (
	"fmt"
	"time"

	"github.com/rs/xid"
)

func NewMeta() *Meta {
	return &Meta{
		Id:       xid.New().String(),
		CreateAt: time.Now().Unix(),
	}
}

func (m *Meta) IdWithPrefix(prefix string) *Meta {
	m.Id = fmt.Sprintf("%s-%s", prefix, m.Id)
	return m
}
