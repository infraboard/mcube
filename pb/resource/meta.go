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
		Label:    make(map[string]string),
	}
}

func (m *Meta) IdWithPrefix(prefix string) *Meta {
	m.Id = fmt.Sprintf("%s-%s", prefix, m.Id)
	return m
}

func NewScope() *Scope {
	return &Scope{}
}
