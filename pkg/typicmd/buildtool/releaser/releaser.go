package releaser

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/go-github/v27/github"

	"github.com/typical-go/typical-go/pkg/typictx"
	"github.com/typical-go/typical-go/pkg/typienv"
	"github.com/typical-go/typical-go/pkg/utility/bash"
	"github.com/typical-go/typical-go/pkg/utility/git"
	"golang.org/x/oauth2"
)

// Releaser responsible to release distruction
type Releaser struct {
	*typictx.Context
	Force bool
	Alpha bool
}

// Git release and return change logs
func (r *Releaser) Git() (changeLogs []string, err error) {
	git.Fetch()
	defer git.Fetch()
	status := git.Status()
	if !r.Force && status != "" {
		err = fmt.Errorf("Please commit changes first:\n%s", status)
		return
	}
	latestTag := git.LatestTag()
	if !r.Force && latestTag == r.ReleaseTag(r.Alpha) {
		log.Errorf("%s already released", latestTag)
		return
	}
	changeLogs = git.Logs(latestTag)
	if !r.Force && len(changeLogs) < 1 {
		msg := "No change to be released"
		log.Errorf(msg)
		err = errors.New(msg)
		return
	}
	return
}

// Distribution to release the distribution
func (r *Releaser) Distribution() (binaries []string, err error) {
	mainPackage := typienv.App.SrcPath
	for _, target := range r.ReleaseTargets {
		chunks := strings.Split(target, "/")
		if len(chunks) != 2 {
			err = fmt.Errorf("Invalid target '%s': it should be '$GOOS/$GOARCH'", target)
			return
		}
		binary := r.ReleaseBinary(chunks[0], chunks[1], r.Alpha)
		binaryPath := fmt.Sprintf("%s/%s", typienv.Release, binary)
		log.Infof("Create release binary for %s: %s", target, binaryPath)
		// TODO: support cgo
		envs := []string{
			"GOOS=" + chunks[0],
			"GOARCH=" + chunks[1],
		}
		if err = bash.GoBuild(binaryPath, mainPackage, envs...); err != nil {
			return
		}
		binaries = append(binaries, binary)
	}
	return
}

// ReleaseToGithub to release to Github
func (r *Releaser) ReleaseToGithub(binaries, changeLogs []string) (err error) {
	if r.Github == nil {
		return errors.New("Missing Github field")
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Environment 'GITHUB_TOKEN' is missing")
	}
	owner := r.Github.Owner
	repo := r.Github.RepoName
	ctx0 := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx0, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
	if r.isGithubReleased(ctx0, client.Repositories) {
		msg := fmt.Sprintf("Release for %s/%s (%s) already exist", owner, repo, r.ReleaseTag(r.Alpha))
		log.Info(msg)
		return errors.New(msg)
	}
	log.Info("Generate release note")
	var rn strings.Builder
	for _, changelog := range changeLogs {
		if !ignoring(changelog) {
			rn.WriteString(changelog)
			rn.WriteString("\n")
		}
	}
	log.Infof("Create github release for %s/%s", owner, repo)
	var release *github.RepositoryRelease
	if release, err = r.createGithubRelease(ctx0, client.Repositories, rn.String()); err != nil {
		return
	}
	for _, binary := range binaries {
		log.Infof("Upload asset: %s", binary)
		if err = r.uploadToGithub(ctx0, client.Repositories, *release.ID, binary); err != nil {
			return
		}
	}
	return
}

func (r *Releaser) isGithubReleased(ctx context.Context, service *github.RepositoriesService) bool {
	owner := r.Github.Owner
	repo := r.Github.RepoName
	tag := r.ReleaseTag(r.Alpha)
	_, _, err := service.GetReleaseByTag(ctx, owner, repo, tag)
	return err == nil
}

func (r *Releaser) createGithubRelease(ctx context.Context, service *github.RepositoriesService, releaseNote string) (release *github.RepositoryRelease, err error) {
	releaseTag := r.ReleaseTag(r.Alpha)
	release, _, err = service.CreateRelease(ctx,
		r.Github.Owner,
		r.Github.RepoName,
		&github.RepositoryRelease{
			Name:       github.String(fmt.Sprintf("%s - %s", r.ReleaseName(), releaseTag)),
			TagName:    github.String(releaseTag),
			Body:       github.String(releaseNote),
			Draft:      github.Bool(false),
			Prerelease: github.Bool(r.Alpha),
		},
	)
	return
}

func (r *Releaser) uploadToGithub(ctx context.Context, service *github.RepositoriesService, repoID int64, binary string) (err error) {
	binaryPath := fmt.Sprintf("%s/%s", typienv.Release, binary)
	var file *os.File
	if file, err = os.Open(binaryPath); err != nil {
		return
	}
	_, _, err = service.UploadReleaseAsset(ctx,
		r.Github.Owner,
		r.Github.RepoName,
		repoID,
		&github.UploadOptions{
			Name: binary,
		},
		file,
	)
	return
}

// ReleaseTag to get release tag
func (r *Releaser) ReleaseTag(alpha bool) string {
	var builder strings.Builder
	builder.WriteString("v")
	builder.WriteString(r.Version)
	if alpha {
		builder.WriteString("-alpha")
	}
	if r.WithGitBranch {
		builder.WriteString("_")
		builder.WriteString(git.Branch())
	}
	if r.WithLatestGitCommit {
		builder.WriteString("_")
		builder.WriteString(git.LatestCommit())
	}
	return builder.String()
}

// ReleaseName to get release name
func (r *Releaser) ReleaseName() string {
	name := r.Name
	if name == "" {
		dir, _ := os.Getwd()
		name = filepath.Base(dir)
	}
	return name
}

// ReleaseBinary to get release binary
func (r *Releaser) ReleaseBinary(os1, arch string, alpha bool) string {
	return strings.Join([]string{
		r.ReleaseName(),
		r.ReleaseTag(alpha),
		os1,
		arch,
	}, "_")
}
