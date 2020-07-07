package wrapper

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typlog"
	"github.com/typical-go/typical-go/pkg/typtmpl"
	"github.com/urfave/cli/v2"
)

const (
	projectPkgVar = "github.com/typical-go/typical-go/pkg/typgo.ProjectPkg"
	typicalTmpVar = "github.com/typical-go/typical-go/pkg/typgo.TypicalTmp"
	gitignore     = ".gitignore"
	typicalw      = "typicalw"
)

type (
	wrapper struct {
		Name    string
		Version string
		typlog.Logger
	}
)

func (w *wrapper) app() *cli.App {
	app := cli.NewApp()
	app.Name = w.Name
	app.Usage = ""       // NOTE: intentionally blank
	app.Description = "" // NOTE: intentionally blank
	app.Version = w.Version

	app.Commands = []*cli.Command{
		{
			Name:  "wrap",
			Usage: "wrap the project with its build-tool",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "typical-tmp", Value: ".typical-tmp"},
				&cli.StringFlag{Name: "descriptor-pkg", Value: "typical"},
				&cli.StringFlag{Name: "project-pkg"},
			},
			Action: w.wrap,
		},
	}

	return app
}

func (w *wrapper) wrap(c *cli.Context) error {
	typicalTmp := c.String("typical-tmp")
	descriptorPkg := c.String("descriptor-pkg")

	projectPkg, err := w.projectPkg(c)
	if err != nil {
		return err
	}

	if err := typgo.Wrap(typicalTmp, projectPkg); err != nil {
		return err
	}

	if err := w.generateTypicalwIfNotExist(typicalTmp, projectPkg); err != nil {
		return err
	}

	if err := w.generateBuildMain(projectPkg, descriptorPkg); err != nil {
		return err
	}

	chksum := createChecksum(descriptorPkg)

	if _, err = os.Stat(typgo.BuildToolBin); os.IsNotExist(err) || !isSameChecksum(typgo.BuildChecksum, chksum) {
		if err = saveChecksum(typgo.BuildChecksum, chksum); err != nil {
			return err
		}

		w.Info("Build the build-tool")
		gobuild := &execkit.GoBuild{
			Out:    typgo.BuildToolBin,
			Source: "./" + typgo.BuildToolSrc,
			Ldflags: []string{
				execkit.BuildVar(projectPkgVar, projectPkg),
				execkit.BuildVar(typicalTmpVar, typicalTmp),
			},
		}

		return gobuild.Run(c.Context)
	}
	return nil
}

func (w *wrapper) projectPkg(c *cli.Context) (string, error) {
	projectPkg := c.String("project-pkg")
	if projectPkg != "" {
		return projectPkg, nil
	}
	var stderr strings.Builder
	var stdout strings.Builder
	cmd := execkit.Command{
		Name:   "go",
		Args:   []string{"list", "-m"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if err := cmd.Run(c.Context); err != nil {
		return "", errors.New(stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (w *wrapper) generateTypicalwIfNotExist(typicalTmp, projectPkg string) error {
	if _, err := os.Stat(typicalw); !os.IsNotExist(err) {
		return nil
	}

	w.Infof("Generate %s", typicalw)
	f, err := os.OpenFile("typicalw", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := &typtmpl.Typicalw{
		TypicalSource: "github.com/typical-go/typical-go/cmd/typical-go",
		TypicalTmp:    typicalTmp,
		ProjectPkg:    projectPkg,
	}
	return tmpl.Execute(f)
}

func (w *wrapper) generateBuildMain(projectPkg, descriptorPkg string) error {
	src := typgo.BuildToolSrc + "/main.go"
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		return nil
	}

	f, err := os.OpenFile(src, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := &typtmpl.BuildToolMain{
		DescPkg: fmt.Sprintf("%s/%s", projectPkg, descriptorPkg),
	}
	return tmpl.Execute(f)
}

func createChecksum(source string) []byte {
	h := sha256.New()
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if b, err := ioutil.ReadFile(path); err == nil {
			h.Write(b)
		}
		return nil
	})
	return h.Sum(nil)
}

func isSameChecksum(filename string, checksum []byte) bool {
	if b, err := ioutil.ReadFile(filename); err == nil {
		return bytes.Compare(checksum, b) == 0
	}
	return false
}

func saveChecksum(filename string, checksum []byte) error {
	return ioutil.WriteFile(filename, checksum, 0777)
}
