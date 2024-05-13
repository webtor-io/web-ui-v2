package web

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	csrf "github.com/utrack/gin-csrf"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"

	"github.com/webtor-io/web-ui-v2/services"
)

const (
	webHostFlag       = "host"
	webPortFlag       = "port"
	sessionSecretFlag = "secret"
	assetsPathFlag    = "assets-path"
	assetsHostFlag    = "assets-host"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
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

		cli.StringFlag{
			Name:   assetsPathFlag,
			Usage:  "assets path",
			Value:  "./assets/dist",
			EnvVar: "ASSETS_PATH",
		},
		cli.StringFlag{
			Name:   assetsHostFlag,
			Usage:  "assets host",
			Value:  "",
			EnvVar: "WEB_ASSETS_HOST",
		},
	)
}

type Web struct {
	host string
	port int
	ln   net.Listener
	r    *gin.Engine
}

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

func (s *Web) Close() {
	log.Info("closing Web")
	defer func() {
		log.Info("Web closed")
	}()
	if s.ln != nil {
		s.ln.Close()
	}
}

func New(c *cli.Context, r *gin.Engine) (*Web, error) {
	var (
		store sessions.Store
		err   error
	)

	if c.String(services.RedisHostFlag) != "" && c.Int(services.RedisPortFlag) != 0 {
		url := fmt.Sprintf("%v:%v", c.String(services.RedisHostFlag), c.Int(services.RedisPortFlag))
		store, err = redis.NewStore(10, "tcp", url, "", []byte(sessionSecretFlag))
		if err != nil {
			return nil, err
		}
		log.Infof("using redis store %v", url)
	} else {
		store = cookie.NewStore([]byte(sessionSecretFlag))
	}
	r.Use(sessions.Sessions("session", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: c.String(sessionSecretFlag),
		ErrorFunc: func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/auth/dashboard") {
				c.Next()
				return
			}
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	r.UseRawPath = true
	assetsPath := c.String(assetsPathFlag)
	r.Static("/assets", assetsPath)
	r.Static("/pub", "./pub")
	r.StaticFile("/favicon.ico", assetsPath+"/favicon.ico")

	return &Web{
		host: c.String(webHostFlag),
		port: c.Int(webPortFlag),
		r:    r,
	}, nil
}
