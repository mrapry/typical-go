package typical

import (
	"github.com/typical-go/typical-go/examples/generate-mock/internal/helloworld"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
)

// Descriptor of sample
var Descriptor = typgo.Descriptor{
	Name:    "generate-mock",
	Version: "1.0.0",

	EntryPoint: helloworld.Main,

	Build: &typgo.StdBuild{},

	Utility: typmock.Utility(),

	Layouts: []string{
		"internal",
	},
}
