package golang

import (
	"io"
	"os"
)

// SourceCode is source code recipe for generated.go in typical package
type SourceCode struct {
	Imports

	PackageName string
	Structs     []Struct
	Init        *Function
}

// NewSourceCode return new instance of SourceCode
func NewSourceCode(pkgName string) *SourceCode {
	return &SourceCode{
		PackageName: pkgName,
		Init:        &Function{Name: "init"},
	}
}

func (r SourceCode) Write(w io.Writer) (err error) {
	writelnf(w, "// Autogenerated by Typical-Go. DO NOT EDIT.\n")
	writelnf(w, "package %s", r.PackageName)
	for _, importPogo := range r.Imports {
		writelnf(w, `import %s "%s"`, importPogo.Name, importPogo.Path)
	}
	for i := range r.Structs {
		r.Structs[i].Write(w)
	}
	if !r.Init.IsEmpty() {
		r.Init.Write(w)
	}
	return
}

// WriteToFile to write to file
func (r SourceCode) WriteToFile(filename string) (err error) {
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return
	}
	defer f.Close()
	return r.Write(f)
}

// AddStruct to add struct
func (r *SourceCode) AddStruct(structs ...Struct) *SourceCode {
	r.Structs = append(r.Structs, structs...)
	return r
}