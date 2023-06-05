package handlers

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mpaas/swagger"
)

// API Doc
func APIDocs(apiDocPath string) *restful.WebService {
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       apiDocPath,
		PostBuildSwaggerObjectHandler: swagger.Docs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}

	return restfulspec.NewOpenAPIService(config)
}
