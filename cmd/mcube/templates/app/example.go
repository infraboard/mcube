package app

// ExampleHTTPObjTemplate todo
const ExampleHTTPObjTemplate = `package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/label"

	"{{.PKG}}/client"
	"{{.PKG}}/apps"
	"{{.PKG}}/apps/example"
)

var (
	api = &handler{}
)

type handler struct {
	service example.ServiceClient
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("examples")

	r.BasePath("books")
	r.Handle("POST", "/", h.CreateBook).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryBook).AddLabel(label.Get)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Example()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("example", api)
}
`

// ExampleHTTPMethodTemplate todo
const ExampleHTTPMethodTemplate = `package http

import (
	"net/http"

	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/http/request"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"{{.PKG}}/apps/example"
)

func (h *handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := &example.CreateBookRequest{}

	var header, trailer metadata.MD
	ins, err := h.service.CreateBook(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryBook(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := example.NewQueryBookRequest(page)

	var header, trailer metadata.MD
	dommains, err := h.service.QueryBook(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, dommains)
}
`

// ExamplePBRequestTemplate todo
const ExamplePBRequestTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/apps/example";

import "github.com/infraboard/mcube/pb/page/page.proto";

// CreateBookRequest 创建Book请求
message CreateBookRequest {
    // book名称
    string name = 1;
}

// QueryBookRequest 查询Book请求
message QueryBookRequest {
    infraboard.mcube.page.PageRequest page = 1;
    string name = 2;
}
`

// ExamplePBResponseTemplate todo
const ExamplePBResponseTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/apps/example";

// Book todo
message Book {
    // 唯一ID
    string id = 1;
    // 名称
    string name =2;
}

// BookSet todo
message BookSet {
    int64 total = 1;
    repeated Book items = 2;
}
`

// ExamplePBServiceTemplate todo
const ExamplePBServiceTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/apps/example";

import "apps/example/pb/request.proto";
import "apps/example/pb/reponse.proto";

service Service {
	rpc CreateBook(CreateBookRequest) returns(Book);
	rpc QueryBook(QueryBookRequest) returns(BookSet);
}
`

// ExampleIMPLOBJTemplate todo
const ExampleIMPLOBJTemplate = `package impl

import (

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/pb/http"

	"{{.PKG}}/apps"
	"{{.PKG}}/apps/example"
)

var (
	// Service 服务实例
	Service = &service{}
)

type service struct {
	example.UnimplementedServiceServer

	log logger.Logger
}

func (s *service) Config() error {
    // get global config with here
	s.log = zap.L().Named("Example")
	return nil
}

// HttpEntry todo
func (s *service) HTTPEntry() *http.EntrySet {
	return example.HttpEntry()
}

func init() {
	pkg.RegistryService("example", Service)
}
`

// ExampleIMPLMethodTemplate todo
const ExampleIMPLMethodTemplate = `package impl

import (
	"context"

	"github.com/infraboard/mcube/grpc/gcontext"

	"{{.PKG}}/apps"
	"{{.PKG}}/apps/example"
)

func (s *service) CreateBook(ctx context.Context, req *example.CreateBookRequest) (*example.Book, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := pkg.S().GetToken(in.GetRequestID())
	s.log.Debug(tk)
	return example.NewBook(req), nil
}

func (s *service) QueryBook(ctx context.Context, req *example.QueryBookRequest) (*example.BookSet, error) {
	return example.NewBookSet(), nil
}
`

// ExampleRequestExtTemplate 模板
const ExampleRequestExtTemplate = `package example

import "github.com/infraboard/mcube/http/request"

// NewQueryBookRequest 查询book列表
func NewQueryBookRequest(page *request.PageRequest) *QueryBookRequest {
	return &QueryBookRequest{
		Page: &page.PageRequest,
	}
}
`

// ExampleResponseExtTemplate 模板
const ExampleResponseExtTemplate = `package example

// NewBook todo
func NewBook(req *CreateBookRequest) *Book {
	return &Book{
		Id:   "mock id",
		Name: req.Name,
	}
}

// NewBookSet 实例
func NewBookSet() *BookSet {
	return &BookSet{
		Items: []*Book{},
	}
}
`
