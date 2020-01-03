package typbuildtool

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/typical-go/typical-go/pkg/typenv"
	"github.com/urfave/cli/v2"
)

// Run the build tool
func Run(d *typcore.ProjectDescriptor) {
	var filenames []string
	var dirs []string
	var err error
	if dirs, filenames, err = projectFiles(typenv.Layout.App); err != nil {
		log.Fatal(err.Error())
	}
	buildtool := buildtool{
		ProjectDescriptor: d,
		filenames:         filenames,
		dirs:              dirs,
	}
	app := cli.NewApp()
	app.Name = d.Name
	app.Usage = ""
	app.Description = d.Description
	app.Version = d.Version
	app.Before = buildtool.before
	app.Commands = buildtool.commands(d)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}

type buildtool struct {
	*typcore.ProjectDescriptor
	filenames []string
	dirs      []string
}

func (t buildtool) commands(d *typcore.ProjectDescriptor) (cmds []*cli.Command) {
	cmds = []*cli.Command{
		t.cmdBuild(),
		cmdClean(),
		t.cmdRun(),
		cmdTest(),
		t.cmdMock(),
	}
	if t.Releaser != nil {
		cmds = append(cmds, t.cmdRelease())
	}
	cmds = append(cmds, BuildCommands(t.ProjectDescriptor)...)
	return
}

// BuildCommands return list of command
func BuildCommands(d *typcore.ProjectDescriptor) (cmds []*cli.Command) {
	ctx := typcore.NewContext(d)
	for _, module := range d.AllModule() {
		if commander, ok := module.(typcore.BuildCommander); ok {
			for _, cmd := range commander.BuildCommands(ctx) {
				cmds = append(cmds, cmd)
			}
		}
	}
	return
}
