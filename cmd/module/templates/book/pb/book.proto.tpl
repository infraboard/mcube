syntax = "proto3";

package demo.book;
option go_package = "{{.PKG}}/apps/book";

import "github.com/infraboard/mcube/pb/page/page.proto";
import "github.com/infraboard/mcube/pb/request/request.proto";

service Service {
    rpc CreateBook(CreateBookRequest) returns(Book);
    rpc QueryBook(QueryBookRequest) returns(BookSet);
    rpc DescribeBook(DescribeBookRequest) returns(Book);
    rpc UpdateBook(UpdateBookRequest) returns(Book);
    rpc DeleteBook(DeleteBookRequest) returns(Book);
}

// Book todo
message Book {
    // 唯一ID
    // @gotags: json:"id" bson:"_id"
    string id = 1;
    // 录入时间
    // @gotags: json:"create_at" bson:"create_at"
    int64 create_at = 2;
    // 更新时间
    // @gotags: json:"update_at" bson:"update_at"
    int64 update_at = 3;
    // 更新人
    // @gotags: json:"update_by" bson:"update_by"
    string update_by = 4;
    // 书本信息
    // @gotags: json:"data" bson:"data"
    CreateBookRequest data = 5;
}

message CreateBookRequest {
    // 创建人
    // @gotags: json:"create_by" bson:"create_by"
    string create_by = 1;
    // 名称
    // @gotags: json:"name" bson:"name" validate:"required"
    string name = 2;
    // 作者
    // @gotags: json:"author" bson:"author" validate:"required"
    string author = 3;
}

message QueryBookRequest {
    // 分页参数
    // @gotags: json:"page" 
    infraboard.mcube.page.PageRequest page = 1;
    // 关键字参数
    // @gotags: json:"keywords"
    string keywords = 2;  
}

// BookSet todo
message BookSet {
    // 分页时，返回总数量
    // @gotags: json:"total"
    int64 total = 1;
    // 一页的数据
    // @gotags: json:"items"
    repeated Book items = 2;
}

message DescribeBookRequest {
    // book id
    // @gotags: json:"id"
    string id = 1;
}

message UpdateBookRequest {
    // book id
    // @gotags: json:"id"
    string id = 1;
    // 更新模式
    // @gotags: json:"update_mode"
    infraboard.mcube.request.UpdateMode update_mode = 2;
    // 更新人
    // @gotags: json:"update_by"
    string update_by = 3;
    // 更新时间
    // @gotags: json:"update_at"
    int64 update_at = 4;
    // 更新的书本信息
    // @gotags: json:"data"
    CreateBookRequest data = 5;
}

message DeleteBookRequest {
    // book id
    // @gotags: json:"id"
    string id = 1;
}