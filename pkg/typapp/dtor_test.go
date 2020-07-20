package typapp_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
)

func TestDtorAnnotation_Execute(t *testing.T) {
	target := "some-target"
	defer os.Remove(target)
	dtorAnnot := &typapp.DtorAnnotation{Target: target}
	ctx := &typast.Context{
		Context: &typgo.Context{
			BuildSys: &typgo.BuildSys{
				Descriptor: &typgo.Descriptor{},
			},
		},
		ASTStore: &typast.ASTStore{
			Annots: []*typast.Annot{
				{TagName: "dtor", Decl: &typast.Decl{Name: "Clean", Package: "pkg", Type: typast.FuncType}},
			},
		},
	}

	require.NoError(t, dtorAnnot.Annotate(ctx))

	b, _ := ioutil.ReadFile(target)
	require.Equal(t, []byte(`package main

// Autogenerated by Typical-Go. DO NOT EDIT.

import (
)

func init() { 
	typapp.AppendDtor(
		&typapp.Destructor{Fn: pkg.Clean},
	)
}`), b)

}

func TestDtorAnnotation_GetTarget(t *testing.T) {
	testcases := []struct {
		TestName string
		*typapp.DtorAnnotation
		Context  *typast.Context
		Expected string
	}{
		{
			TestName:       "initial target is not set",
			DtorAnnotation: &typapp.DtorAnnotation{},
			Context: &typast.Context{
				Context: &typgo.Context{
					BuildSys: &typgo.BuildSys{
						Descriptor: &typgo.Descriptor{Name: "name0"},
					},
				},
			},
			Expected: "cmd/name0/dtor_annotated.go",
		},
		{
			TestName: "initial target is set",
			DtorAnnotation: &typapp.DtorAnnotation{
				Target: "some-target",
			},
			Expected: "some-target",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.Expected, tt.GetTarget(tt.Context))
		})
	}
}
