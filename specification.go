package saunter

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	spec "github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

const rootDir = "routers"

var specs = map[string]spec.Swagger{}

type parser struct{ routes gin.RoutesInfo }

func (p *parser) parseAttribute(value string) string {
	return strings.Split(value, " ")[0]
}

func (p *parser) parseInfo(context []string) *spec.Info {
	scheme := &spec.Info{}
	for _, value := range context {
		attribute := p.parseAttribute(value)
		term := strings.TrimSpace(value[len(attribute):])
		switch attribute {
		case "@title":
			scheme.Title = term
		case "@version":
			scheme.Version = term
		case "@description":
			scheme.Description = term
		}
	}
	return scheme
}

func (p *parser) parseSecurity(context []string) *spec.SecuritySchemeRef {
	scheme := &spec.SecuritySchemeRef{
		Value: &spec.SecurityScheme{},
	}
	for _, raw := range context {
		attribute := p.parseAttribute(raw)
		value := strings.TrimSpace(raw[len(attribute):])
		switch attribute {
		case "@name":
			scheme.Value.Name = value
		case "@type":
			scheme.Value.Type = value
		case "@description":
			scheme.Value.Description = value
		case "@scheme":
			scheme.Value.Scheme = value
		}
	}
	return scheme
}

func (p *parser) parseMainFile(swagger *spec.Swagger, astFile *ast.File) {
	securitySchemes := map[string]*spec.SecuritySchemeRef{}
	for _, comment := range astFile.Comments {
		context := strings.Split(comment.Text(), "\n")
		for i, value := range context {
			switch p.parseAttribute(value) {
			case "@Info":
				swagger.Info = p.parseInfo(context[i+1:])
			case "@Security":
				security := p.parseSecurity(context[i+1:])
				securitySchemes[security.Value.Name] = security
			}
		}
	}
	swagger.Components.SecuritySchemes = securitySchemes
}

func (p *parser) parseRouterFile(swagger *spec.Swagger, exactPath string, astFile *ast.File) {
	for _, route := range p.routes {
		if strings.EqualFold(route.Path, exactPath) {
			for _, topDecl := range astFile.Decls {
				switch decl := topDecl.(type) {
				case *ast.FuncDecl:
					if decl.Doc != nil {
						operation := &spec.Operation{
							Parameters: spec.NewParameters(),
							Responses:  make(spec.Responses),
							Security:   spec.NewSecurityRequirements(),
						}
						for _, raw := range strings.Split(decl.Doc.Text(), "\n") {
							attribute := p.parseAttribute(raw)
							value := strings.TrimSpace(raw[len(attribute):])
							switch attribute {
							case "@security":
								operation.Security.With(spec.NewSecurityRequirement().Authenticate(value))
							case "@summary":
								operation.Summary = value
							default:
								if r := regexp.MustCompile(`^@(\d+)`); r.MatchString(attribute) {
									matched := r.FindStringSubmatch(attribute)
									code, _ := strconv.Atoi(matched[1])
									response := spec.NewResponse().WithDescription(value).WithJSONSchema(&spec.Schema{})
									operation.AddResponse(code, response)
								}
							}
						}
						swagger.AddOperation(route.Path, route.Method, operation)
					}
				}
			}
		}
	}
}

func generateSpecs(basePath string, routes gin.RoutesInfo) {
	type file struct {
		name string
		path string
	}
	p := &parser{routes}
	files := map[string][]file{}
	filepath.Walk(rootDir, func(path string, f os.FileInfo, err error) error {
		if r := regexp.MustCompile(`/(v\d+)/(.+).go$`); r.MatchString(path) {
			matched := r.FindStringSubmatch(path)
			version := matched[1]
			name := matched[2]
			files[version] = append(files[version], file{name, path})
		}
		return nil
	})
	for version, content := range files {
		swagger := &spec.Swagger{OpenAPI: "3.0.0"}
		for _, file := range content {
			astFile, err := goparser.ParseFile(token.NewFileSet(), file.path, nil, goparser.ParseComments)
			if err != nil {
				panic(fmt.Errorf("cannot parse file %s: %v", file.path, err))
			}
			switch file.name {
			case "main":
				p.parseMainFile(swagger, astFile)
			default:
				exactPath := strings.Join([]string{basePath, version, file.name}, "/") + "/"
				p.parseRouterFile(swagger, exactPath, astFile)
			}
		}
		specs[version] = *swagger
	}
}

func determineSpec(routePath string) struct{ Spec spec.Swagger } {
	return struct{ Spec spec.Swagger }{specs[strings.Split(routePath, "/")[2]]}
}
