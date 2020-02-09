package stdbuild

import (
	"context"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typenv"

	"github.com/urfave/cli/v2"
)

// CmdBuild is command for build
func CmdBuild(bc *typcore.BuildContext) *cli.Command {
	return &cli.Command{
		Name:    "build",
		Aliases: []string{"b"},
		Usage:   "Build the binary",
		Action: func(c *cli.Context) (err error) {
			log.Info("Build the application")
			return buildProject(c.Context, bc)
		},
	}
}

func buildProject(ctx context.Context, bc *typcore.BuildContext) (err error) {
	if bc.Build != nil {
		if err = bc.Build.Prebuild(ctx, bc); err != nil {
			return
		}
	}
	cmd := exec.CommandContext(ctx,
		"go",
		"build",
		"-o", typenv.AppBin,
		"./"+typenv.AppMainPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}