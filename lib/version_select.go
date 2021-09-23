package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	ErrVersionNotSelected = errors.New("version is not in the constraint range")
)

func KctlVersionList(constraint string, log *zap.SugaredLogger) ([]string, error) {

	c, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Errorf("The constraint \"%s\" is not valid.", constraint)
		return nil, err
	}

	log.Debugf("Found a valid constraint \"%s\"; fetching tags", constraint)

	// TODO: use a local cache for git
	// TODO: manage cache expiration

	// asd, _ := getTemporaryGitDir(log)
	// log.Info(asd)

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/kubernetes/kubernetes.git"},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	var versions []*semver.Version
	// Filters the references list and only keeps tags
	for _, ref := range refs {
		if ref.Name().IsTag() {
			t := ref.Name().Short()
			if err := validateTag(t, *c); err == nil {
				v, _ := semver.NewVersion(t)
				versions = append(versions, v)
			}
		}
	}
	sort.Sort(semver.Collection(versions))
	if len(versions) == 0 {
		log.Error("No version is satisfying the constraint!")
	}

	log.Infof("Tags found: %v", versions)
	tags := []string{}
	for _, v := range versions {
		tags = append(tags, v.String())
	}
	return tags, nil
}

func validateTag(tag string, constraint semver.Constraints) error {
	v, err := semver.NewVersion(tag)
	if err != nil {
		return err
	}
	if !constraint.Check(v) {
		return ErrVersionNotSelected
	}
	return nil
}

func getTemporaryGitDir(log *zap.SugaredLogger) (path string, err error) {
	parentDir := os.TempDir()
	globPattern := filepath.Join(parentDir, "kctlswitch-*")
	matches, _ := filepath.Glob(globPattern)
	var tempGitDir string
	if len(matches) == 0 {
		log.Debug("Didn't find a match for a temp dir, building one.")
		tempGitDir, err = ioutil.TempDir(parentDir, "kctlswitch-*")
		if err != nil {
			return "", err
		}

	} else {
		log.Debugf("Found a matching temp dir \"%s\", using that", matches[0])
		tempGitDir = matches[0]
	}
	tempGitPath := fmt.Sprintf("%s/.git", tempGitDir)
	os.MkdirAll(tempGitPath, 0755)
	return tempGitPath, nil
}
