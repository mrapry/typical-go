package app

const ctxSrc = `package typical

import (
	"{{.Pkg}}/app"

	"github.com/typical-go/typical-go/pkg/typcore"
)

// Context of Project
var Context = &typcore.Context{
	Name:      "{{.Name}}",
	Version:   "0.0.1",
	Package:   "{{.Pkg}}",
	AppModule: &app.Module{},
}
`

const blankCtxSrc = `package typical

import (
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Context of Project
var Context = &typcore.Context{
	Name:      "{{.Name}}",
	Version:   "0.0.1",
	Package:   "{{.Pkg}}",
}
`
