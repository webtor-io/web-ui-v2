package web

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/gin-contrib/multitemplate"
	sv "github.com/webtor-io/web-ui-v2/services"
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
	assetsPath string
}

func NewTemplateHandler(c *cli.Context) *TemplateHandler {
	return &TemplateHandler{
		assetsHost: c.String(assetsHostFlag),
		assetsPath: c.String(sv.AssetsPathFlag),
	}
}

func (s *TemplateHandler) MakeAsset(in string) string {
	h, _ := s.getAssetHash(in)
	return s.assetsHost + "/assets/" + in + "?" + h
}

func (s *TemplateHandler) getAssetHash(name string) (string, error) {
	f, err := os.Open(s.assetsPath + "/" + name)
	if err != nil {
		return "", err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *TemplateHandler) RegisterTemplate(r multitemplate.Renderer, name string, layouts []string, partials []string, fm template.FuncMap) {
	funcs := template.FuncMap{
		"asset":         s.MakeAsset,
		"makeJobLogURL": MakeJobLogURL,
		"bitsForHumans": BitsForHumans,
		"log":           Log,
		"shortErr":      ShortErr,
	}
	for k, v := range fm {
		funcs[k] = v
	}
	pp := []string{}
	for _, p := range partials {
		pp = append(pp, fmt.Sprintf("templates/partials/%v.html", p))
	}
	for _, l := range layouts {
		templates := []string{}
		templates = append(templates, fmt.Sprintf("templates/layouts/%v.html", l), fmt.Sprintf("templates/%v.html", name))
		templates = append(templates, pp...)
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
