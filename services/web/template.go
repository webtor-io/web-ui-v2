package web

import (
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/gin-contrib/multitemplate"
)

const (
	assetsHostFlag = "assets-host"
)

func RegisterTemplateHandlerFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   assetsHostFlag,
			Usage:  "assets host",
			Value:  "",
			EnvVar: "WEB_ASSETS_HOST",
		},
	)
}

type TemplateHandler struct {
	assetsHost string
}

func NewTemplateHandler(c *cli.Context) *TemplateHandler {
	return &TemplateHandler{
		assetsHost: c.String(assetsHostFlag),
	}
}

func (s *TemplateHandler) RegisterTemplate(r multitemplate.Renderer, name string, layouts []string) {
	funcs := template.FuncMap{
		"asset":           func(in string) interface{} { return s.assetsHost + "/assets/" + in },
		"makeBreadcrumbs": MakeBreadcrumbs,
		"hasBreadcrumbs":  HasBreadcrumbs,
		"hasPagination":   HasPagination,
		"makePagination":  MakePagination,
		"makeJobLogURL":   MakeJobLogURL,
		"bitsForHumans":   BitsForHumans,
		"log":             Log,
		"shortErr":        ShortErr,
	}
	partials, _ := filepath.Glob("templates/partials/*.html")
	for _, l := range layouts {
		templates := []string{}
		templates = append(templates, fmt.Sprintf("templates/layouts/%v.html", l), fmt.Sprintf("templates/%v.html", name))
		templates = append(templates, partials...)
		r.AddFromFilesFuncs(fmt.Sprintf("%v_%v", name, l), funcs, templates...)
	}
}

type Template struct {
	name string
	c    *gin.Context
	d    any
}

func (s *Template) R(code int) {
	name := s.name + "_standard"
	if s.c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		if s.c.GetHeader("X-Layout") != "" {
			name = s.name + "_" + s.c.GetHeader("X-Layout")
		} else {
			name = s.name + "_async"
		}
		s.c.Header("X-Template", name)
	}
	s.c.HTML(code, name, s.d)
}

func (s *TemplateHandler) MakeTemplate(c *gin.Context, name string, d any) *Template {
	return &Template{
		name: name,
		c:    c,
		d:    d,
	}
}
