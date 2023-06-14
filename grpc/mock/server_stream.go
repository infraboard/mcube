package mock

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func NewServerStreamBase() *ServerStreamBase {
	return &ServerStreamBase{}
}

type ServerStreamBase struct {
}

func (i *ServerStreamBase) SetHeader(metadata.MD) error {
	return nil
}
func (i *ServerStreamBase) SendHeader(metadata.MD) error {
	return nil
}
func (i *ServerStreamBase) SetTrailer(metadata.MD) {

}
func (i *ServerStreamBase) Context() context.Context {
	return context.Background()
}
func (i *ServerStreamBase) SendMsg(m interface{}) error {
	return nil
}
func (i *ServerStreamBase) RecvMsg(m interface{}) error {
	return nil
}
