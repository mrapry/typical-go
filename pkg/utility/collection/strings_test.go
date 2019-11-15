package collection_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/utility/collection"
)

func TestStrings(t *testing.T) {
	var coll collection.Strings
	coll.Add("hello")
	coll.Add("world")
	require.EqualValues(t, []string{"hello", "world"}, coll)
}
