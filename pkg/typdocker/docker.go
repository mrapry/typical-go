package typdocker

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-go/pkg/typbuildtool"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

// Docker of docker
type Docker struct {
	version   Version
	composers []Composer
}

// Composer responsible to compose docker
type Composer interface {
	DockerCompose(version Version) *ComposeObject
}

// New docker module
func New() *Docker {
	return &Docker{
		version: "3",
	}
}

// WithVersion to set the version
func (m *Docker) WithVersion(version Version) *Docker {
	m.version = version
	return m
}

// WithComposers to set the composers
func (m *Docker) WithComposers(composers ...Composer) *Docker {
	m.composers = composers
	return m
}

// Commands of docker
func (m *Docker) Commands(ctx *typbuildtool.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "docker",
			Usage: "Docker utility",
			Subcommands: []*cli.Command{
				{
					Name:  "compose",
					Usage: "Generate docker-compose.yaml",
					Action: func(c *cli.Context) (err error) {
						if len(m.composers) < 1 {
							return errors.New("No composers is set")
						}
						var out []byte
						log.Info("Generate docker-compose.yml")
						if out, err = yaml.Marshal(m.dockerCompose()); err != nil {
							return
						}
						if err = ioutil.WriteFile("docker-compose.yml", out, 0644); err != nil {
							return
						}
						return
					},
				},
				{
					Name:  "up",
					Usage: "Spin up docker containers according docker-compose",
					Action: func(ctx *cli.Context) (err error) {
						cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						return cmd.Run()
					},
				},
				{
					Name:  "down",
					Usage: "Take down all docker containers according docker-compose",
					Action: func(ctx *cli.Context) (err error) {
						cmd := exec.Command("docker-compose", "down")
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						return cmd.Run()
					},
				},
				{
					Name:  "wipe",
					Usage: "Kill all running docker container",
					Action: func(ctx *cli.Context) (err error) {
						var builder strings.Builder
						cmd := exec.Command("docker", "ps", "-q")
						cmd.Stderr = os.Stderr
						cmd.Stdout = &builder
						if err = cmd.Run(); err != nil {
							return
						}
						dockerIDs := strings.Split(builder.String(), "\n")
						for _, id := range dockerIDs {
							if id != "" {
								cmd := exec.Command("docker", "kill", id)
								cmd.Stderr = os.Stderr
								if err = cmd.Run(); err != nil {
									return
								}
							}
						}
						return
					},
				},
			},
		},
	}
}

func (m *Docker) dockerCompose() (root *ComposeObject) {
	root = &ComposeObject{
		Version:  m.version,
		Services: make(Services),
		Networks: make(Networks),
		Volumes:  make(Volumes),
	}
	for _, composer := range m.composers {
		if obj := composer.DockerCompose(m.version); obj != nil {
			root.Append(obj)
		}
	}
	return
}
