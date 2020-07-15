package typgo

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typtmpl"
)

var (
	dtorTag = "dtor"
)

type (
	// DtorAnnotation represent @dtor annotation
	DtorAnnotation struct {
		Target string
	}
)

var _ Action = (*DtorAnnotation)(nil)

// Execute @dtor
func (a *DtorAnnotation) Execute(c *Context) error {
	var dtors []*typtmpl.Dtor
	for _, annot := range c.ASTStore.Annots {
		if annot.Check(dtorTag, typast.FuncType) {
			dtors = append(dtors, typtmpl.CreateDtor(annot))
		}
	}
	return writeGoSource(
		a.GetTarget(c),
		&typtmpl.DtorAnnotated{
			Package: "main",
			Imports: c.Imports,
			Dtors:   dtors,
		},
	)
}

// GetTarget to get generation target for dtor
func (a *DtorAnnotation) GetTarget(c *Context) string {
	if a.Target == "" {
		a.Target = fmt.Sprintf("cmd/%s/dtor_annotated.go", c.Descriptor.Name)
	}
	return a.Target
}
