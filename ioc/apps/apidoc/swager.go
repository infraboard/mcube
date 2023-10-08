package apidoc

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

// API Doc
func APIDocs(apiDocPath string, docs restfulspec.PostBuildSwaggerObjectFunc) *restful.WebService {
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       apiDocPath,
		PostBuildSwaggerObjectHandler: docs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}

	return restfulspec.NewOpenAPIService(config)
}
