package typfactory_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typfactory"
)

type testcase struct {
	testName string
	typfactory.Writer
	expected string
}

func TestTypicalw(t *testing.T) {

	testcases := []testcase{
		{
			Writer: &typfactory.Typicalw{
				TypicalSource: "some-source",
				TypicalTmp:    "some-tmp",
			},
			expected: `#!/bin/bash

set -e

TYPSRC=some-source
TYPTMP=some-tmp
TYPGO=$TYPTMP/bin/typical-go

if ! [ -s $TYPGO ]; then
	go build -o $TYPGO $TYPSRC
fi

$TYPGO wrap \
	-typical-tmp=$TYPTMP \

$TYPTMP/bin/build-tool $@
`,
		},
		{
			Writer: &typfactory.Typicalw{
				TypicalSource: "some-source",
				TypicalTmp:    "some-tmp",
				ProjectPkg:    "some-project-pkg",
			},
			expected: `#!/bin/bash

set -e

TYPSRC=some-source
TYPTMP=some-tmp
TYPGO=$TYPTMP/bin/typical-go

if ! [ -s $TYPGO ]; then
	go build -o $TYPGO $TYPSRC
fi

$TYPGO wrap \
	-typical-tmp=$TYPTMP \
	-project-pkg="some-project-pkg" \

$TYPTMP/bin/build-tool $@
`,
		},
	}

	for _, tt := range testcases {
		var debugger strings.Builder
		require.NoError(t, tt.Write(&debugger))
		require.Equal(t, tt.expected, debugger.String())
	}

}
