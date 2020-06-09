package typmock_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typmock"
)

var (
	target1 = &typmock.Mock{Pkg: "pkg1", Dir: "dir1", Source: "target1"}
	target2 = &typmock.Mock{Pkg: "pkg1", Dir: "dir1", Source: "target2"}
	target3 = &typmock.Mock{Pkg: "pkg2", Dir: "dir2", Source: "target3"}
	target4 = &typmock.Mock{Pkg: "pkg1", Dir: "dir1", Source: "target4"}
	target5 = &typmock.Mock{Pkg: "pkg1", Dir: "dir1", Source: "target5"}
	target6 = &typmock.Mock{Pkg: "pkg2", Dir: "dir2", Source: "target6"}

	targets = []*typmock.Mock{
		target1,
		target2,
		target3,
		target4,
		target5,
		target6,
	}

	dir1 = []*typmock.Mock{
		target1,
		target2,
		target4,
		target5,
	}

	dir2 = []*typmock.Mock{
		target3,
		target6,
	}
)

func TestTargetMap(t *testing.T) {
	m := typmock.TargetMap{}

	for _, mock := range targets {
		m.Put(mock)
	}

	require.Equal(t, typmock.TargetMap{"dir1": dir1, "dir2": dir2}, m.Filter("dir1", "dir2"))
	require.Equal(t, typmock.TargetMap{"dir1": dir1}, m.Filter("dir1"))
	require.Equal(t, typmock.TargetMap{}, m.Filter("not-found"))
}