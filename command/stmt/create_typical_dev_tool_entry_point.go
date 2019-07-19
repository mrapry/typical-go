package stmt

import (
	"html/template"
	"os"

	"github.com/typical-go/typical-go/typicore"
)

const typicalDevToolEntryPointTemplate = `package main

import (
	"log"
	"os"

	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typimain"
	"{{ .ProjectPath }}/typical"
)

func main() {
	cli := typimain.NewTypicalTaskTool(typical.Context)
	err := cli.Cli().Run(os.Args)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
`

type CreateTypicalDevToolEntryPoint struct {
	Metadata *typicore.ContextMetadata
	Target   string
}

func (c CreateTypicalDevToolEntryPoint) Run() (err error) {
	f, err := os.Create(c.Target)
	if err != nil {
		return
	}

	tmpl, _ := template.New("typical_context").Parse(typicalDevToolEntryPointTemplate)
	return tmpl.Execute(f, c.Metadata)
}
