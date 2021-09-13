package lib

import (
	"log"

	"github.com/go-git/go-git/v5" // with go modules enabled (GO111MODULE=on or outside GOPATH)
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/Masterminds/semver"
)

type kctlVersionList struct {
	kctllist []string
}

func BuildKctlList(kubectlGitRepo string) kctlVersionList {
	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{kubectlGitRepo},
	})

	log.Print("Fetching tags...")

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Filters the references list and only keeps tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			_, err := semver.NewVersion(string(ref.Name().Short()))
			if err == nil {
				tags = append(tags, ref.Name().Short())
			}
		}
	}

	if len(tags) == 0 {
		log.Println("No tags!")
	}

	log.Printf("Tags found: %v", tags)

	return kctlVersionList{
		kctllist: tags,
	}
}

/*
func SelectKctlVersion(versionConstraint string, versionList kctlVersionList, includeUnstable bool) {
	_, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		print("Nope")
	}
}
*/
