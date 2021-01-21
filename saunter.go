package saunter

import (
	"html/template"
	"net/http"

	spec "github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
)

// Saunter implements Swagger specs generator
type Saunter struct {
	template *template.Template
	specs    map[string]spec.Swagger
}

// Initialize generate Saunter required data
func Initialize(basePath string, routes gin.RoutesInfo) {
	generateIndexTemplate()
	generateSpecs(basePath, routes)
}

// Handler wraps `http.Handler` into `gin.HandlerFunc`
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := indexTemplate.Execute(c.Writer, determineSpec(c.FullPath())); err != nil {
			panic(err)
		}
	}
}

// Static creates Saunter static file system
func Static() http.FileSystem {
	static, err := fs.New()
	if err != nil {
		panic(err)
	}

	return static
}
