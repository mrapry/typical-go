package typbuildtool

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typenv"

	"github.com/urfave/cli/v2"
)

func (t buildtool) cmdBuild() *cli.Command {
	return &cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "Build the binary",
		Action:  t.buildBinary,
	}
}

func (t buildtool) buildBinary(c *cli.Context) (err error) {
	if err = t.prebuild(c); err != nil {
		return
	}
	log.Info("Build the application")
	cmd := exec.CommandContext(c.Context,
		"go",
		"build",
		"-o", typenv.AppBin,
		"./"+typenv.AppMainPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
