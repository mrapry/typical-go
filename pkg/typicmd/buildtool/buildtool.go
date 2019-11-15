package buildtool

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typicli"
	"github.com/typical-go/typical-go/pkg/utility/bash"

	"github.com/typical-go/typical-go/pkg/typicmd/buildtool/releaser"
	"github.com/typical-go/typical-go/pkg/typictx"
	"github.com/typical-go/typical-go/pkg/typienv"
	"github.com/urfave/cli"
)

type buildtool struct {
	*typictx.Context
}

func (t buildtool) commands() (cmds []cli.Command) {
	cmds = []cli.Command{
		{
			Name:      "build",
			ShortName: "b",
			Usage:     "Build the binary",
			Action:    t.buildBinary,
		},
		{
			Name:      "clean",
			ShortName: "c",
			Usage:     "Clean the project from generated file during build time",
			Action:    t.cleanProject,
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the binary",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-build",
					Usage: "Run the binary without build",
				},
			},
			Action: t.runBinary,
		},
		{
			Name:      "test",
			ShortName: "t",
			Usage:     "Run the testing",
			Action:    t.runTesting,
		},
		{
			Name:  "release",
			Usage: "Release the distribution",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-test",
					Usage: "Release without run automated test",
				},
				cli.BoolFlag{
					Name:  "no-github",
					Usage: "Release without create github release",
				},
				cli.BoolFlag{
					Name:  "force",
					Usage: "Release by passed all validation",
				},
				cli.BoolFlag{
					Name:  "alpha",
					Usage: "Release for alpha version",
				},
			},
			Action: t.releaseDistribution,
		},
		{
			Name:  "mock",
			Usage: "Generate mock class",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-delete",
					Usage: "Generate mock class with delete previous generation",
				},
			},
			Action: t.generateMock,
		},
	}
	if t.ReadmeGenerator != nil {
		cmds = append(cmds, cli.Command{
			Name:  "readme",
			Usage: "Generate readme document",
			Action: func(*cli.Context) error {
				return t.ReadmeGenerator.GenerateReadme(t.Context)
			},
		})
	}
	cmds = append(cmds, Commands(t.Context)...)
	return
}

func (t buildtool) cliBefore(ctx *cli.Context) (err error) {
	return t.Context.Validate()
}

func (t buildtool) buildBinary(ctx *cli.Context) error {
	log.Info("Build the application")
	return bash.GoBuild(typienv.App.BinPath, typienv.App.SrcPath)
}

func (t buildtool) cleanProject(ctx *cli.Context) (err error) {
	log.Info("Clean the application")
	log.Infof("  Remove %s", typienv.Bin)
	if err = os.RemoveAll(typienv.Bin); err != nil {
		return
	}
	log.Infof("  Remove %s", typienv.Metadata)
	if err = os.RemoveAll(typienv.Metadata); err != nil {
		return
	}
	log.Info("  Remove .env")
	os.Remove(".env")
	return filepath.Walk(typienv.Dependency.SrcPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			log.Infof("  Remove %s", path)
			return os.Remove(path)
		}
		return nil
	})
}

func (t buildtool) runBinary(ctx *cli.Context) (err error) {
	if !ctx.Bool("no-build") {
		if err = t.buildBinary(ctx); err != nil {
			return
		}
	}
	log.Info("Run the application")
	return bash.Run(typienv.App.BinPath, []string(ctx.Args())...)
}

func (t buildtool) runTesting(ctx *cli.Context) error {
	log.Info("Run testings")
	return bash.GoTest(t.TestTargets)
}

func (t buildtool) generateMock(ctx *cli.Context) (err error) {
	log.Info("Generate mocks")
	if err = bash.GoGet("github.com/golang/mock/mockgen"); err != nil {
		return
	}
	mockPkg := typienv.Mock
	if !ctx.Bool("no-delete") {
		log.Infof("Clean mock package '%s'", mockPkg)
		os.RemoveAll(mockPkg)
	}
	for _, mockTarget := range t.MockTargets {
		dest := mockPkg + "/" + mockTarget[strings.LastIndex(mockTarget, "/")+1:]
		err = bash.RunGoBin("mockgen",
			"-source", mockTarget,
			"-destination", dest,
			"-package", mockPkg)
	}
	return
}

func (t buildtool) releaseDistribution(ctx *cli.Context) (err error) {
	log.Info("Release distribution")
	var binaries, changeLogs []string
	if !ctx.Bool("no-test") {
		if err = t.runTesting(ctx); err != nil {
			return
		}
	}
	rel := releaser.Releaser{
		Release: t.Release,
		Force:   ctx.Bool("force"),
		Alpha:   ctx.Bool("alpha"),
	}
	if changeLogs, err = rel.Git(); err != nil {
		return
	}
	if binaries, err = rel.Distribution(); err != nil {
		return
	}
	if !ctx.Bool("no-github") {
		if err = rel.ReleaseToGithub(binaries, changeLogs); err != nil {
			return
		}
	}
	return
}

// Commands return list of command
func Commands(ctx *typictx.Context) (cmds []cli.Command) {
	for _, module := range ctx.AllModule() {
		if cmd := Command(ctx, module); cmd != nil {
			cmds = append(cmds, *cmd)
		}
	}
	return
}

// Command of module
func Command(ctx *typictx.Context, module interface{}) *cli.Command {
	if commander, ok := module.(typicli.BuildCommander); ok {
		cmd := commander.BuildCommand(&typicli.ContextCli{
			Context: ctx,
		})
		return &cmd
	}
	if commander, ok := module.(typicli.Commander); ok {
		cmd := commander.Command(&typicli.Cli{
			Obj: module,
		})
		return &cmd
	}
	return nil
}
