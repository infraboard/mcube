package pkg

// ExampleHTTPTemplate todo
const ExampleHTTPTemplate = `package all

import (
	// 加载服务模块
)
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
