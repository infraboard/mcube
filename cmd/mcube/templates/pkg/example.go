package pkg

// ExampleHTTPObjTemplate todo
const ExampleHTTPObjTemplate = `package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"

	"{{.PKG}}/client"
	"{{.PKG}}/pkg"
	"{{.PKG}}/pkg/example"
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
	r.Handle("POST", "/", h.CreateBook)
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

	httpcontext "github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"{{.PKG}}/pkg"
	"{{.PKG}}/pkg/example"
)

func (h *handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, err := pkg.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := example.NewQueryDomainRequest(page)

	var header, trailer metadata.MD
	ins, err := h.service.CreateBook(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, pkg.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, ins)
	return
}
`

// ExamplePBRequestTemplate todo
const ExamplePBRequestTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/pkg/example";

// CreateBookRequest 创建Book请求
message CreateBookRequest {
    // 应用名称
    string name = 1;
}
`

// ExamplePBResponseTemplate todo
const ExamplePBResponseTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/pkg/example";

import "github.com/infraboard/protoc-gen-go-ext/extension/tag/tag.proto";

// Book todo
message Book {
    // 唯一ID
    string id = 1[
        (google.protobuf.field_tag) = {struct_tag: 'bson:"_id" json:"id"'}
        ];
    // 名称
    string name =2[
        (google.protobuf.field_tag) = {struct_tag: 'bson:"name" json:"name"'}
        ];
}

// BookSet todo
message BookSet {
    int64 total = 1[
        (google.protobuf.field_tag) = {struct_tag: 'bson:"total" json:"total"'}
        ];
    repeated Book items = 2[
        (google.protobuf.field_tag) = {struct_tag: 'bson:"items" json:"items"'}
        ];
}
`

// ExamplePBServiceTemplate todo
const ExamplePBServiceTemplate = `syntax = "proto3";

package {{.Name}}.example;
option go_package = "{{.PKG}}/pkg/example";

import "pkg/example/pb/request.proto";
import "pkg/example/pb/reponse.proto";
import "github.com/infraboard/mcube/pb/http/entry.proto";

service Service {
	rpc CreateBook(CreateBookRequest) returns(Book) {
		option (mcube.http.rest_api) = {
			path: "/applications/"
			method: "POST"
			resource: "application"
			auth_enable: true
			permission_enable: true
			labels: [{
				key: "action"
				value: "create"
			}]
		};
	};
}
`

// ExampleIMPLOBJTemplate todo
const ExampleIMPLOBJTemplate = `package impl

import (

	"github.com/infraboard/mcube/pb/http"

	"{{.PKG}}/pkg"
	"{{.PKG}}/pkg/example"
)

var (
	// Service 服务实例
	Service = &service{}
)

type service struct {
	example.UnimplementedServiceServer
}

func (s *service) Config() error {
    // get global config with here

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

	"{{.PKG}}/pkg/example"
)

func (s *service) CreateBook(ctx context.Context, req *example.CreateBookRequest) (*example.Book, error) {
	return &example.Book{}, nil
}`
