package static

import (
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

const (
	AssetsPathFlag = "assets-path"
	AssetsHostFlag = "assets-host"
)

func RegisterFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   AssetsPathFlag,
			Usage:  "assets path",
			Value:  "./assets/dist",
			EnvVar: "ASSETS_PATH",
		},
		cli.StringFlag{
			Name:   AssetsHostFlag,
			Usage:  "assets host",
			Value:  "",
			EnvVar: "WEB_ASSETS_HOST",
		},
	)
}

func RegisterHandler(c *cli.Context, r *gin.Engine) {
	assetsPath := c.String(AssetsPathFlag)
	r.Static("/assets", assetsPath)
	r.Static("/pub", "./pub")
	r.StaticFile("/favicon.ico", assetsPath+"/favicon.ico")
}
