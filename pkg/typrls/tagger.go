package typrls

import (
	"fmt"

	"github.com/typical-go/typical-go/pkg/typgo"
)

type (
	// Tagger responsible to create tag
	Tagger interface {
		CreateTag(c *typgo.Context, alpha bool) string
	}
	// StdTagger standard tagger
	StdTagger struct {
		Extra string // Add extra information after TagName e.g. v0.0.1+extra
		GitID bool   // Add latest git id after TagName e.g. v0.0.1+5339d71
	}
)

//
// StdTagger
//

var _ Tagger = (*StdTagger)(nil)

// WithExtra with extra information
func (s *StdTagger) WithExtra(extra string) *StdTagger {
	s.Extra = extra
	return s
}

// WithGitID with latest git id as extra information
func (s *StdTagger) WithGitID() *StdTagger {
	s.GitID = true
	return s
}

// CreateTag create tag
func (s *StdTagger) CreateTag(c *typgo.Context, alpha bool) string {
	tagName := "v0.0.1"
	if c.BuildSys.ProjectVersion != "" {
		tagName = fmt.Sprintf("v%s", c.BuildSys.ProjectVersion)
	}
	if s.Extra != "" {
		tagName = tagName + "+" + s.Extra
	}
	if s.GitID {
		latestGitID := latestGitID(c.Ctx())
		if len(latestGitID) > 6 {
			latestGitID = latestGitID[:6]
		}
		tagName = tagName + "+" + latestGitID
	}
	if alpha {
		tagName = tagName + "_alpha"
	}

	return tagName
}
