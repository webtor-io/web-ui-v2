package services

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	csrf "github.com/utrack/gin-csrf"
	"net"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

const (
	webHostFlag       = "host"
	webPortFlag       = "port"
	sessionSecretFlag = "secret"
)

func RegisterWebFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   webHostFlag,
			Usage:  "listening host",
			Value:  "",
			EnvVar: "WEB_HOST",
		},
		cli.IntFlag{
			Name:   webPortFlag,
			Usage:  "http listening port",
			Value:  8080,
			EnvVar: "WEB_PORT",
		},
		cli.StringFlag{
			Name:   sessionSecretFlag,
			Usage:  "session secret",
			Value:  "secret123",
			EnvVar: "SESSION_SECRET",
		},
	)
}

type Web struct {
	host   string
	port   int
	secret string
	ln     net.Listener
	r      *gin.Engine
	re     multitemplate.Renderer
}

const (
	assetsPath string = "./assets/dist"
)

func (s *Web) Serve() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	ln, err := net.Listen("tcp", addr)
	s.ln = ln
	if err != nil {
		return errors.Wrap(err, "failed to web listen to tcp connection")
	}
	log.Infof("serving Web at %v", addr)
	return http.Serve(s.ln, s.r)
}

type Handler interface {
	RegisterRoutes(r *gin.Engine)
	RegisterTemplates(r multitemplate.Renderer)
}

func (s *Web) RegisterHandler(h Handler) {
	h.RegisterRoutes(s.r)
	h.RegisterTemplates(s.re)
}

func (s *Web) Close() {
	log.Info("closing Web")
	defer func() {
		log.Info("Web closed")
	}()
	if s.ln != nil {
		s.ln.Close()
	}
}

func NewWeb(c *cli.Context) *Web {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: c.String(sessionSecretFlag),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	r.UseRawPath = true
	r.Static("/assets", assetsPath)
	r.StaticFile("/favicon.ico", assetsPath+"/favicon.ico")
	re := multitemplate.NewRenderer()
	r.HTMLRender = re

	return &Web{
		host: c.String(webHostFlag),
		port: c.Int(webPortFlag),
		r:    r,
		re:   re,
	}
}
