package typgo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/typical-go/typical-go/pkg/git"
	"github.com/urfave/cli/v2"
)

func cmdPublish(c *BuildTool) *cli.Command {
	return &cli.Command{
		Name:    "publish",
		Usage:   "Publish the project",
		Aliases: []string{"p"},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "no-test", Usage: "skip the test"},
			&cli.BoolFlag{Name: "force", Usage: "Release by passed all validation"},
			&cli.BoolFlag{Name: "alpha", Usage: "Release for alpha version"},
		},
		Action: c.ActionFunc("PUBLISH", Publish),
	}
}

// Publish project
func Publish(c *CliContext) (err error) {

	var (
		releaseFiles []string
		latest       string
		gitLogs      []*git.Log
	)

	if !c.Cli.Bool("no-test") {
		if err = test(c); err != nil {
			return
		}
	}

	ctx := c.Cli.Context

	if err = git.Fetch(ctx); err != nil {
		return fmt.Errorf("Failed git fetch: %w", err)
	}
	defer git.Fetch(ctx)

	force := c.Cli.Bool("force")

	if status := git.Status(ctx); status != "" && !force {
		return fmt.Errorf("Please commit changes first:\n%s", status)
	}

	alpha := c.Cli.Bool("alpha")
	tag := releaseTag(c, alpha)

	if latest = git.LatestTag(ctx); latest == tag && !force {
		return fmt.Errorf("%s already released", latest)
	}

	if gitLogs = git.RetrieveLogs(ctx, latest); len(gitLogs) < 1 && !force {
		return errors.New("No change to be released")
	}

	rc := &ReleaseContext{
		CliContext: c,
		Alpha:      alpha,
		Tag:        tag,
		GitLogs:    gitLogs,
	}

	if releaseFiles, err = release(rc); err != nil {
		return
	}

	if err = publish(&PublishContext{
		ReleaseContext: rc,
		ReleaseFiles:   releaseFiles,
	}); err != nil {
		return
	}

	return

}

// Publish the project
func publish(c *PublishContext) (err error) {
	for _, module := range c.BuildTool.BuildSequences {
		if publisher, ok := module.(Publisher); ok {
			if err = publisher.Publish(c); err != nil {
				return
			}
		}
	}
	return
}

func release(c *ReleaseContext) (files []string, err error) {
	for _, module := range c.BuildTool.BuildSequences {
		if releaser, ok := module.(Releaser); ok {
			var files1 []string
			if files1, err = releaser.Release(c); err != nil {
				return
			}
			files = append(files, files1...)
		}
	}
	return
}

func releaseTag(c *CliContext, alpha bool) string {
	var builder strings.Builder
	builder.WriteString("v")
	builder.WriteString(c.BuildTool.Version)
	// if c.BuildTool.IncludeBranch {
	// 	builder.WriteString("_")
	// 	builder.WriteString(git.Branch(c.Context))
	// }
	// if c.BuildTool.IncludeCommitID {
	// 	builder.WriteString("_")
	// 	builder.WriteString(git.LatestCommit(c.Context))
	// }
	if alpha {
		builder.WriteString("_alpha")
	}
	return builder.String()
}