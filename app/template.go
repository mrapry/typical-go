package app

const appSrc = `package app

import (
	"fmt"
	"{{.Pkg}}/app/config"

	"github.com/typical-go/typical-go/pkg/typcfg"
)

// Module of application
type Module struct {}

// Action of application
func (*Module) Action() interface{} {
	return func(cfg config.Config) {
		fmt.Printf("Hello %s\n", cfg.Hello)
	}
}

// Configure the application
func (*Module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "APP"
	spec = &config.Config{}
	loadFn = func(loader typcfg.Loader) (cfg config.Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}
`

const appSrcTest = `package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typmodule"

	"{{.Pkg}}/app"
)

func TestModule(t *testing.T) {
	a := &app.Module{}
	require.True(t, typmodule.IsActionable(a))
	require.True(t, typcfg.IsConfigurer(a))
}
`

const configSrc = "package config\n\n// Config of app\ntype Config struct {\n	Hello string `default:\"World\"`\n}"

const ctxSrc = `package typical

import (
	"{{.Pkg}}/app"

	"github.com/typical-go/typical-go/pkg/typctx"
)

// Context of Project
var Context = &typctx.Context{
	Name:      "{{.Name}}",
	Version:   "0.0.1",
	Package:   "{{.Pkg}}",
	AppModule: &app.Module{},
}
`

const blankCtxSrc = `package typical

import (
	"github.com/typical-go/typical-go/pkg/typctx"
)

// Context of Project
var Context = &typctx.Context{
	Name:      "{{.Name}}",
	Version:   "0.0.1",
	Package:   "{{.Pkg}}",
}
`

const typicalw = `#!/bin/bash
set -e

CHECKSUM_DATA=$(cksum {{.ContextFile}})

if ! [ -s {{.ChecksumFile}} ]; then
	mkdir -p {{.LayoutMetadata}}
	cksum typical/context.go > {{.ChecksumFile}}
else
	CHECKSUM_UPDATED=$([ "$CHECKSUM_DATA" == "$(cat {{.ChecksumFile}} )" ] ; echo $?)
fi

if [ "$CHECKSUM_UPDATED" == "1" ] || ! [[ -f {{.PrebuilderBin}} ]] ; then 
	echo $CHECKSUM_DATA > {{.ChecksumFile}}
	echo "Build the prebuilder"
	go build -o {{.PrebuilderBin}} ./{{.PrebuilderMainPath}}
fi

./{{.PrebuilderBin}} $CHECKSUM_UPDATED
./{{.BuildtoolBin}} $@`

const gomod = `module {{.Pkg}}

go 1.13

require github.com/typical-go/typical-go v{{.TypicalVersion}}
`

const gitignore = `/bin
/release
/.typical-metadata
/vendor
.envrc
.env
*.test
*.out`

const moduleSrc = `package {{.Name}}

import (
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/urfave/cli/v2"
)

// Config of {{.Name}}
type Config struct {
	// TODO:
}

// Module of {{.Name}}
type Module struct {}

// Configure the module
func (m *Module) Configure() (prefix string, spec, loadFn interface{}) {
	prefix = "{{.Prefix}}"
	spec = &Config{}
	loadFn = func(loader typcfg.Loader) (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide the dependencies
func (m *Module) Provide() []interface{} {
	return []interface{}{
		// TODO: (1) put functions to be provided as dependencies
		// TODO: (2) remove this function if not required
	}
}

// Prepare the module
func (m *Module) Prepare() []interface{} {
	return []interface{}{
		// TODO: (1) put functions that run before the application
		// TODO: (2) remove this function if not required
	}
}

// Destroy the dependencies
func (m *Module) Destroy() []interface{} {
	return []interface{}{
		// TODO: (1) functions to destroy dependencies
		// TODO: (2) remove this function if not required
	}
}

// BuildCommands is commands to exectuce from Build-Tool
func (m *Module) BuildCommands(c *typcli.ModuleCli) []*cli.Command {
	return []*cli.Command{
		// TODO: (1) add command to execute from Build-Tool
		// TODO: (2) remove this function if not required
	}
}

`

const moduleSrcTest = `package {{.Name}}

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/typical-go/typical-go/pkg/typcli"
	"github.com/typical-go/typical-go/pkg/typmodule"
)

func TestModule(t *testing.T) {
	m := &Module{}
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typmodule.IsDestroyer(m))
	require.True(t, typmodule.IsProvider(m))
	require.True(t, typcfg.IsConfigurer(m))
	require.True(t, typcli.IsBuildCommander(m))
}
`