package static

import (
	"os"
	"path/filepath"
	"strings"

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

func RegisterHandler(c *cli.Context, r *gin.Engine) error {
	assetsPath := c.String(AssetsPathFlag)
	pubPath := "pub"

	r.Static("/assets", assetsPath)
	r.Static("/pub", pubPath)

	err := filepath.Walk(pubPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			r.StaticFile(strings.TrimPrefix(path, pubPath), path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	r.StaticFile("/favicon.ico", assetsPath+"/favicon.ico")
	return nil
}
