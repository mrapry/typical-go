package typmod_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typmod"
)

type SampleAttribute struct {
	Name        string
	Description string
}

func TestName(t *testing.T) {
	testcases := []struct {
		obj  interface{}
		name string
	}{
		{
			obj: SampleAttribute{
				Name: "some-name",
			},
			name: "some-name",
		},
		{
			obj:  SampleAttribute{},
			name: "SampleAttribute",
		},
		{
			obj:  struct{}{},
			name: "struct {}",
		},
		{
			obj:  nil,
			name: "nil",
		},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.name, typmod.Name(tt.obj))
	}
}

func TestDescription(t *testing.T) {
	testcases := []struct {
		obj         interface{}
		description string
	}{
		{
			obj: SampleAttribute{
				Description: "some-description",
			},
			description: "some-description",
		},
		{
			obj:         SampleAttribute{},
			description: "",
		},
		{
			obj:         struct{}{},
			description: "",
		},
		{
			obj:         nil,
			description: "nil",
		},
	}
	for _, tt := range testcases {
		require.Equal(t, tt.description, typmod.Description(tt.obj))
	}
}
