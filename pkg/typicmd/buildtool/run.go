package buildtool

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typictx"
	"github.com/urfave/cli"
)

// Run the build tool
func Run(c *typictx.Context) {
	buildtool := buildtool{Context: c}
	app := cli.NewApp()
	app.Name = c.Name
	app.Usage = ""
	app.Description = c.Description
	app.Version = c.Release.Version
	app.Before = buildtool.cliBefore
	app.Commands = buildtool.commands()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
