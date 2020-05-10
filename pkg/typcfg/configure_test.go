package typcfg_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

func TestConfigurers(t *testing.T) {
	cfg1 := &typcfg.Configuration{}
	cfg2 := &typcfg.Configuration{}
	cfg3 := &typcfg.Configuration{}
	cfg4 := &typcfg.Configuration{}
	configurers := typcfg.Configurers{
		cfg1,
		cfg2,
		typcfg.Configurers{
			cfg3,
			cfg4,
		},
	}
	require.Equal(t, []*typcfg.Configuration{
		cfg1,
		cfg2,
		cfg3,
		cfg4,
	}, configurers.Configurations())

}