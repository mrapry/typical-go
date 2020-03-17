package typicalgo

import (
	"os"

	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// TypicalGo is app of typical-go
type TypicalGo struct {
	wrapper typcore.Wrapper
}

// New instance of TypicalGo
func New() *TypicalGo {
	return &TypicalGo{
		wrapper: typcore.NewWrapper(),
	}
}

// RunApp to run the typical-go
func (t *TypicalGo) RunApp(d *typcore.Descriptor) (err error) {
	app := cli.NewApp()
	app.Name = d.Name
	app.Usage = "" // NOTE: intentionally blank
	app.Description = d.Description
	app.Version = d.Version

	app.Commands = []*cli.Command{
		{
			Name:  "wrap",
			Usage: "wrap the project with its build-tool",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "tmp-folder", Required: true},
				&cli.StringFlag{Name: "project-package", Usage: "To override generated ProjectPackage in context"},
			},
			Action: func(cliCtx *cli.Context) (err error) {
				return t.Wrap(&typcore.WrapContext{
					Descriptor:     d,
					Ctx:            cliCtx.Context,
					TmpFolder:      cliCtx.String("tmp-folder"),
					ProjectPackage: cliCtx.String("project-package"),
				})
			},
		},
	}
	return app.Run(os.Args)
}

// Wrap the project
func (t *TypicalGo) Wrap(c *typcore.WrapContext) error {
	return t.wrapper.Wrap(c)
}