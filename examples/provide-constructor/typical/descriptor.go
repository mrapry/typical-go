package typical

import (
	"github.com/typical-go/typical-go/examples/provide-constructor/helloworld"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuildtool"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Modules
var (
	hello = helloworld.New()
)

// Descriptor of sample
var Descriptor = typcore.Descriptor{
	Name:    "provide-constructor",
	Version: "0.0.1",

	App: typapp.New(hello),

	BuildTool: typbuildtool.New(),
}