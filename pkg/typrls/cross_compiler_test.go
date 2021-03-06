package typrls_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/execkit"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/urfave/cli/v2"
)

func TestCrossCompile(t *testing.T) {
	testcases := []struct {
		TestName string
		typrls.CrossCompiler
		Context         *typrls.Context
		RunExpectations []*execkit.RunExpectation
		ExpectedErr     string
	}{
		{
			CrossCompiler: typrls.CrossCompiler{
				Targets: []typrls.Target{"darwin/amd64", "linux/amd64"},
			},
			Context: &typrls.Context{
				TagName: "v0.0.1",
				Context: &typgo.Context{
					BuildSys: &typgo.BuildSys{
						Descriptor: &typgo.Descriptor{
							ProjectName: "myproject",
						},
					},
					Context: &cli.Context{},
				},
			},
			RunExpectations: []*execkit.RunExpectation{
				{CommandLine: "go build -ldflags \"-X github.com/typical-go/typical-go/pkg/typapp.Name=myproject -X github.com/typical-go/typical-go/pkg/typapp.Version=v0.0.1\" -o /myproject_v0.0.1_darwin_amd64 ./cmd/myproject"},
				{CommandLine: "go build -ldflags \"-X github.com/typical-go/typical-go/pkg/typapp.Name=myproject -X github.com/typical-go/typical-go/pkg/typapp.Version=v0.0.1\" -o /myproject_v0.0.1_linux_amd64 ./cmd/myproject"},
			},
		},
		{
			TestName: "go build error",
			CrossCompiler: typrls.CrossCompiler{
				Targets: []typrls.Target{"darwin/amd64"},
			},
			Context: &typrls.Context{
				TagName: "v0.0.1",
				Context: &typgo.Context{
					BuildSys: &typgo.BuildSys{
						Descriptor: &typgo.Descriptor{
							ProjectName: "myproject",
						},
					},
					Context: &cli.Context{},
				},
			},
			RunExpectations: []*execkit.RunExpectation{
				{
					CommandLine: "go build -ldflags \"-X github.com/typical-go/typical-go/pkg/typapp.Name=myproject -X github.com/typical-go/typical-go/pkg/typapp.Version=v0.0.1\" -o /myproject_v0.0.1_darwin_amd64 ./cmd/myproject",
					ReturnError: errors.New("some-error"),
				},
			},
			ExpectedErr: "some-error",
		},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			unpatch := execkit.Patch(tt.RunExpectations)
			defer unpatch(t)
			err := tt.Release(tt.Context)
			if tt.ExpectedErr != "" {
				require.EqualError(t, err, tt.ExpectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTarget(t *testing.T) {
	testcases := []struct {
		TestName string
		typrls.Target
		ExpectedOS   string
		ExpectedArch string
	}{
		{Target: "darwin/amd64", ExpectedOS: "darwin", ExpectedArch: "amd64"},
		{Target: "linux/amd64", ExpectedOS: "linux", ExpectedArch: "amd64"},
		{Target: "no-slash", ExpectedOS: "", ExpectedArch: ""},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			require.Equal(t, tt.ExpectedOS, tt.OS())
			require.Equal(t, tt.ExpectedArch, tt.Arch())
		})
	}
}
