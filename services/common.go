package services

import (
	"regexp"

	"github.com/urfave/cli"
)

var SHA1R = regexp.MustCompile("(?i)[0-9a-f]{5,40}")

var (
	AssetsPathFlag = "assets-path"
	DomainFlag     = "domain"
	SMTPHostFlag   = "smtp-host"
	SMTPUserFlag   = "smtp-user"
	SMTPPassFlag   = "smtp-pass"
	SMTPPortFlag   = "smtp-port"
	SMTPSecureFlag = "smtp-secure"
	UseAuthFlag    = "use-auth"
)

func RegisterCommonFlags(f []cli.Flag) []cli.Flag {
	f = append(f,
		cli.StringFlag{
			Name:   AssetsPathFlag,
			Usage:  "assets path",
			Value:  "./assets/dist",
			EnvVar: "ASSETS_PATH",
		},
	)
	f = append(f,
		cli.StringFlag{
			Name:   DomainFlag,
			Usage:  "domain",
			Value:  "http://localhost:8080",
			EnvVar: "DOMAIN",
		},
	)
	f = append(f,
		cli.StringFlag{
			Name:   SMTPHostFlag,
			Usage:  "smtp host",
			EnvVar: "SMTP_HOST",
		},
	)
	f = append(f,
		cli.StringFlag{
			Name:   SMTPUserFlag,
			Usage:  "smtp user",
			EnvVar: "SMTP_USER",
		},
	)
	f = append(f,
		cli.StringFlag{
			Name:   SMTPPassFlag,
			Usage:  "smtp pass",
			EnvVar: "SMTP_PASS",
		},
	)
	f = append(f,
		cli.IntFlag{
			Name:   SMTPPortFlag,
			Usage:  "smtp port",
			EnvVar: "SMTP_PORT",
			Value:  465,
		},
	)
	f = append(f,
		cli.BoolTFlag{
			Name:   SMTPSecureFlag,
			Usage:  "smtp secure",
			EnvVar: "SMTP_SECURE",
		},
	)
	f = append(f,
		cli.BoolFlag{
			Name:   UseAuthFlag,
			Usage:  "use auth",
			EnvVar: "USE_AUTH",
		},
	)
	return f
}
