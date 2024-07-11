package main

import (
	"github.com/urfave/cli"
	services "github.com/webtor-io/common-services"
)

func configure(app *cli.App) {
	serveCmd := makeServeCMD()
	migrationCMD := services.MakePGMigrationCMD()
	app.Commands = []cli.Command{serveCmd, migrationCMD}
}
