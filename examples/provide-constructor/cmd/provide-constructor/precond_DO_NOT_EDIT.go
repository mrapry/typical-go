package main

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
	"github.com/typical-go/typical-go/examples/provide-constructor/internal/helloworld"
	"github.com/typical-go/typical-go/pkg/typgo"
)

func init() {
	typgo.Provide(
		&typgo.Constructor{Name: "", Fn: helloworld.HelloWorld},
		&typgo.Constructor{Name: "typical", Fn: helloworld.HelloTypical},
	)
	typgo.Destroy(
		&typgo.Destructor{Fn: helloworld.Close},
	)
}
