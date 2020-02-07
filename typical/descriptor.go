package typical

import (
	"github.com/typical-go/typical-go/app"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typbuild/stdrelease"
	"github.com/typical-go/typical-go/pkg/typcore"
)

// Descriptor of typical-go
var Descriptor = typcore.Descriptor{
	Name:    "Typical-Go",
	Version: app.Version,
	Package: "github.com/typical-go/typical-go",

	App: typapp.New(application),

	Build: typcore.NewBuild().
		WithRelease(stdrelease.New().
			WithPublisher(
				stdrelease.GithubPublisher("typical-go", "typical-go"),
			),
		),
}

var (
	application = app.New()
)
