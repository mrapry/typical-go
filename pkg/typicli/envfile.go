package typicli

import (
	"github.com/typical-go/typical-go/pkg/utility/envfile"
	"github.com/urfave/cli"
)

// LoadEnvFile is cli version of LoadEnvFile
func LoadEnvFile(ctx *cli.Context) (err error) {
	return envfile.Load()
}
