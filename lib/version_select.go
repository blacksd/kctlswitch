package lib

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/semver"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/google/go-github/v39/github"
)

type kctlVersionList struct {
	kctllist []string
}

func BuildKctlList() kctlVersionList {
	client := github.NewClient(nil)

	var tag_list []string

	listOptions := &github.ListOptions{Page: 2, PerPage: 30}
	for {
		tags, response, err := client.Repositories.ListTags(context.TODO(), "kubernetes", "kubectl", listOptions)

		if err != nil {
			log.Fatal(err)
		}

		for _, tag := range tags {
			_, err := semver.NewVersion(tag.GetName())
			if err == nil {
				tag_list = append(tag_list, *tag.Name)
				fmt.Println(*tag.Name)
			}

		}

		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage
	}

	// Filters the references list and only keeps tags
	//
	// for _, release := range releases {
	// 	if release.GetName() {
	// 		_, err := semver.NewVersion(string(ref.Name().Short()))
	// 		if err == nil {
	// 			tags = append(tags, ref.Name().Short())
	// 		}
	// 	}
	// }

	//if len(tags) == 0 {
	//	log.Println("No tags!")
	//}
	//
	//log.Printf("Tags found: %v", tags)

	return kctlVersionList{
		kctllist: tag_list,
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

func FetchGitTags() []string {
	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/kubernetes/kubectl.git"},
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
			// _, err := semver.NewVersion(string(ref.Name().Short()))
			if err := validation.Validate(ref.Name().Short(), validation.By(validateTag)); err == nil {
				tags = append(tags, ref.Name().Short())
			}
		}
	}

	if len(tags) == 0 {
		log.Println("No tags!")
	}

	log.Printf("Tags found: %v", tags)

	return tags
}

// TODO: add constraints as param
func validateTag(value interface{}) error {
	s, _ := value.(string)
	if _, err := semver.NewVersion(s); err != nil {
		return errors.New("The tag '" + s + "' is not semver compliant.")
	}
	return nil
}

// Output: must be abc
