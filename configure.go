package main

import (
	"github.com/urfave/cli"
)

func configure(app *cli.App) {
	serveCmd := makeServeCMD()
	app.Commands = []cli.Command{serveCmd}
}
