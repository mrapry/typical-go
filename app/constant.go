package app

const appSrc = `package app

import "fmt"

// Module of application
func Module() interface{} {
	return &module{}
}

type module struct{}

func (module) Action() interface{} {
	return func() {
		fmt.Println("Hello World")
	}
}
`

const appSrcTest = `package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typmodule"

	"{{.Pkg}}/app"
)

func TestModule(t *testing.T) {
	a := app.Module()
	require.True(t, typmodule.IsActionable(a))
}
`

const ctxSrc = `package typical

import (
	"{{.Pkg}}/app"
	"github.com/typical-go/typical-go/pkg/typctx"
	"github.com/typical-go/typical-go/pkg/typrls"
)

// Context of Project
var Context = &typctx.Context{
	Name:    "{{.Name}}",
	Version: "0.0.1",
	Package: "{{.Pkg}}",
	AppModule: app.Module(),
	Releaser: typrls.Releaser{
		Targets: []typrls.Target{"linux/amd64", "darwin/amd64"},
	},
}
`

const typicalw = `#!/bin/bash
set -e

BIN=${TYPICAL_BIN:-bin}
CMD=${TYPICAL_CMD:-cmd}
BUILD_TOOL=${TYPICAL_BUILD_TOOL:-build-tool}
PRE_BUILDER=${TYPICAL_PRE_BUILDER:-pre-builder}
METADATA=${TYPICAL_METADATA:-.typical-metadata}

CHECKSUM_PATH="$METADATA/checksum "
CHECKSUM_DATA=$(cksum typical/context.go)

if ! [ -s .typical-metadata/checksum ]; then
	mkdir -p $METADATA
	cksum typical/context.go > $CHECKSUM_PATH
else
	CHECKSUM_UPDATED=$([ "$CHECKSUM_DATA" == "$(cat $CHECKSUM_PATH )" ] ; echo $?)
fi

if [ "$CHECKSUM_UPDATED" == "1" ] || ! [[ -f $BIN/$PRE_BUILDER ]] ; then 
	echo $CHECKSUM_DATA > $CHECKSUM_PATH
	echo "Build the pre-builder"
	go build -o $BIN/$PRE_BUILDER ./$CMD/$PRE_BUILDER
fi

./$BIN/$PRE_BUILDER $CHECKSUM_UPDATED
./$BIN/$BUILD_TOOL $@`

const gomod = `module {{.Pkg}}

go 1.12

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