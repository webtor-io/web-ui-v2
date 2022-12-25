package services

import (
	"github.com/urfave/cli"
	"regexp"
)

var SHA1R = regexp.MustCompile("(?i)[0-9a-f]{5,40}")

var AssetsPathFlag = "assets-path"

func RegisterCommonFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   AssetsPathFlag,
			Usage:  "assets path",
			Value:  "./assets/dist",
			EnvVar: "ASSETS_PATH",
		},
	)
}
