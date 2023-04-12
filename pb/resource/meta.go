package resource

import (
	"time"

	"github.com/rs/xid"
)

func NewMeta() *Meta {
	return &Meta{
		Id:       xid.New().String(),
		CreateAt: time.Now().Unix(),
	}
}
