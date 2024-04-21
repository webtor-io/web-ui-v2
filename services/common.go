package services

import (
	"regexp"

	"github.com/urfave/cli"
)

var SHA1R = regexp.MustCompile("(?i)[0-9a-f]{5,40}")

var (
	DomainFlag     = "domain"
	SMTPHostFlag   = "smtp-host"
	SMTPUserFlag   = "smtp-user"
	SMTPPassFlag   = "smtp-pass"
	SMTPPortFlag   = "smtp-port"
	SMTPSecureFlag = "smtp-secure"
	RedisHostFlag  = "redis-host"
	RedisPortFlag  = "redis-port"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {

	f = append(f,
		cli.StringFlag{
			Name:   DomainFlag,
			Usage:  "domain",
			Value:  "http://localhost:8080",
			EnvVar: "DOMAIN",
		},
		cli.StringFlag{
			Name:   SMTPHostFlag,
			Usage:  "smtp host",
			EnvVar: "SMTP_HOST",
		},
		cli.StringFlag{
			Name:   SMTPUserFlag,
			Usage:  "smtp user",
			EnvVar: "SMTP_USER",
		},
		cli.StringFlag{
			Name:   SMTPPassFlag,
			Usage:  "smtp pass",
			EnvVar: "SMTP_PASS",
		},
		cli.IntFlag{
			Name:   SMTPPortFlag,
			Usage:  "smtp port",
			EnvVar: "SMTP_PORT",
			Value:  465,
		},
		cli.BoolTFlag{
			Name:   SMTPSecureFlag,
			Usage:  "smtp secure",
			EnvVar: "SMTP_SECURE",
		},
		cli.StringFlag{
			Name:   RedisHostFlag,
			Usage:  "redis host",
			EnvVar: "REDIS_MASTER_SERVICE_HOST, REDIS_SERVICE_HOST",
		},
		cli.IntFlag{
			Name:   RedisPortFlag,
			Usage:  "redis port",
			EnvVar: "REDIS_MASTER_SERVICE_PORT, REDIS_SERVICE_PORT",
		},
	)

	return f
}
