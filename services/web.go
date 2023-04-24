package services

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	csrf "github.com/utrack/gin-csrf"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

const (
	webHostFlag       = "host"
	webPortFlag       = "port"
	sessionSecretFlag = "secret"
	redisHostFlag     = "redis-host"
	redisPortFlag     = "redis-port"
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
		cli.StringFlag{
			Name:   redisHostFlag,
			Usage:  "redis host",
			EnvVar: "REDIS_MASTER_SERVICE_HOST, REDIS_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   redisPortFlag,
			Usage:  "redis port",
			EnvVar: "REDIS_MASTER_SERVICE_PORT, REDIS_SERVICE_PORT",
		},
	)
}

type Web struct {
	host   string
	port   int
	secret string
	ln     net.Listener
	r      *gin.Engine
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

func NewWeb(c *cli.Context, r *gin.Engine) *Web {
	var store sessions.Store
	if c.String(redisHostFlag) != "" && c.Int(redisPortFlag) != 0 {
		url := fmt.Sprintf("%v:%v", c.String(redisHostFlag), c.Int(redisPortFlag))
		store, _ = redis.NewStore(10, "tcp", url, "", []byte(sessionSecretFlag))
		log.Infof("using redis store %v", url)
	} else {
		store = cookie.NewStore([]byte(sessionSecretFlag))
	}
	r.Use(sessions.Sessions("session", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: c.String(sessionSecretFlag),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
	r.UseRawPath = true
	assetsPath := c.String(AssetsPathFlag)
	r.Static("/assets", assetsPath)
	r.StaticFile("/favicon.ico", assetsPath+"/favicon.ico")

	return &Web{
		host: c.String(webHostFlag),
		port: c.Int(webPortFlag),
		r:    r,
	}
}
