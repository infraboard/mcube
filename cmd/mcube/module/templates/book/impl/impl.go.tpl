package impl

import (
{{ if $.EnableMySQL -}}
	"database/sql"
{{- end }}

{{ if $.EnableMongoDB -}}
	"go.mongodb.org/mongo-driver/mongo"
{{- end }}
	

	"github.com/infraboard/mcube/v2/app"
	"github.com/infraboard/mcube/v2/logger"
	"github.com/infraboard/mcube/v2/logger/zap"
	"google.golang.org/grpc"

	"{{.PKG}}/apps/book"
	"{{.PKG}}/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
{{ if $.EnableMySQL -}}
	db   *sql.DB
{{- end }}
{{ if $.EnableMongoDB -}}
	col *mongo.Collection
{{- end }}
	log  log.Logger
	book.UnimplementedServiceServer
}

func (s *service) Config() error {
{{ if $.EnableMySQL -}}
	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}
	s.db = db
{{- end }}
{{ if $.EnableMongoDB -}}
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}
	s.col = db.Collection(s.Name())
{{- end }}

	s.log = log.Sub(s.Name())
	return nil
}

func (s *service) Name() string {
	return book.AppName
}

func (s *service) Registry(server *grpc.Server) {
	book.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}