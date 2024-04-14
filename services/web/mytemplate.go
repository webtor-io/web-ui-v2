package web

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io"
	"os"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	csrf "github.com/utrack/gin-csrf"
	"github.com/webtor-io/web-ui-v2/services"
)

func NewMyTemplateManager(c *cli.Context, re multitemplate.Renderer) *TemplateManager {
	h := &Helper{
		assetsHost: c.String(assetsHostFlag),
		assetsPath: c.String("assets-path"),
		useAuth:    c.Bool("use-auth"),
		domain:     c.String("domain"),
	}
	return NewTemplateManager(re, template.FuncMap{
		"asset":         h.MakeAsset,
		"devAsset":      h.MakeDevAsset,
		"makeJobLogURL": MakeJobLogURL,
		"bitsForHumans": BitsForHumans,
		"log":           Log,
		"shortErr":      ShortErr,
		"dev":           Dev,
		"useAuth":       h.UseAuth,
		"domain":        h.Domain,
		"has":           Has,
	}, NewContext)
}

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

type Helper struct {
	assetsHost string
	assetsPath string
	useAuth    bool
	domain     string
}

func (s *Helper) UseAuth() bool {
	return s.useAuth
}

func (s *Helper) Domain() string {
	return s.domain
}

func (s *Helper) MakeAsset(in string) string {
	path := s.assetsHost + "/assets/" + in
	if !Dev() {
		h, _ := s.getAssetHash(in)
		path += "?" + h
	}
	return path
}

func (s *Helper) MakeDevAsset(in string) string {
	return s.assetsHost + "/assets/dev/" + in
}

func (s *Helper) getAssetHash(name string) (string, error) {
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

type Context struct {
	Data any
	CSRF string
	Err  error
	User *services.User
}

func NewContext(c *gin.Context, obj any, err error) any {
	u := services.GetUserFromContext(c)
	return &Context{
		Data: obj,
		CSRF: csrf.GetToken(c),
		Err:  err,
		User: u,
	}
}
