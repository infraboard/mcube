package swagger

import (
	"github.com/go-openapi/spec"
	"{{.PKG}}/version"
)

func Docs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "BookService",
			Description: "Resource for managing Books",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "john",
					Email: "john@doe.rp",
					URL:   "http://johndoe.org",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: version.Short(),
		},
	}
	swo.Tags = []spec.Tag{
		{TagProps: spec.TagProps{
			Name:        "books",
			Description: "Managing books"},
		},
	}
}