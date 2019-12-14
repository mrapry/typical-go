package typbuildtool

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typenv"
	"github.com/urfave/cli/v2"
)

func (t buildtool) cmdClean() *cli.Command {
	return &cli.Command{
		Name:    "clean",
		Aliases: []string{"c"},
		Usage:   "Clean the project from generated file during build time",
		Action:  t.cleanProject,
		Subcommands: []*cli.Command{
			{Name: "build-tool", Usage: "Remove build-tool", Action: t.removeBuildTool},
			{Name: "app", Usage: "Remove app", Action: t.removeApp},
			{Name: "metadata", Usage: "Remove metadata", Action: t.removeMetadata},
			{Name: "env", Usage: "Remove envfile", Action: t.removeEnvFile},
			{Name: "constructor", Usage: "Remove constructor", Action: t.removeConstructor},
		},
	}
}

func (t buildtool) cleanProject(ctx *cli.Context) error {
	t.removeBuildTool(ctx)
	t.removeApp(ctx)
	t.removeMetadata(ctx)
	t.removeEnvFile(ctx)
	t.removeConstructor(ctx)
	return nil
}

func (t buildtool) removeBuildTool(ctx *cli.Context) error {
	removeFile(typenv.BuildToolBin)
	return nil
}

func (t buildtool) removeApp(ctx *cli.Context) error {
	removeFile(typenv.AppBin)
	return nil
}

func (t buildtool) removeMetadata(ctx *cli.Context) error {
	removeAllFile(typenv.Layout.Metadata)
	return nil
}

func (t buildtool) removeEnvFile(ctx *cli.Context) error {
	removeFile(".env")
	return nil
}

func (t buildtool) removeConstructor(ctx *cli.Context) error {
	removeFile(typenv.AppMainPath + "/constructor.go")
	return nil
}

func removeFile(name string) {
	log.Infof("Remove: %s", name)
	if err := os.Remove(name); err != nil {
		log.Error(err.Error())
	}
}

func removeAllFile(path string) {
	log.Infof("Remove All: %s", path)
	if err := os.RemoveAll(path); err != nil {
		log.Error(err.Error())
	}
}
