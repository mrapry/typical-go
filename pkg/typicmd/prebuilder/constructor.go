package prebuilder

import (
	"github.com/typical-go/typical-go/pkg/utility/bash"
	"github.com/typical-go/typical-go/pkg/utility/debugkit"

	"github.com/typical-go/typical-go/pkg/typicmd/prebuilder/golang"
	"github.com/typical-go/typical-go/pkg/typienv"
)

type constructor struct {
	ApplicationImports golang.Imports
	Constructors       []string
}

func (g constructor) generate(target string) (err error) {
	defer debugkit.ElapsedTime("Generate constructor")()
	src := golang.NewSourceCode(typienv.Dependency.Package)
	src.Imports = g.ApplicationImports
	src.AddConstructors(g.Constructors...)
	if err = src.Cook(target); err != nil {
		return
	}
	return bash.GoImports(target)
}
