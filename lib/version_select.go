package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v39/github"
)

type kctlVersionList struct {
	kctllist []string
}

func BuildKctlList() kctlVersionList {
	client := github.NewClient(nil)

	listOptions := &github.ListOptions{PerPage: 50}
	for {
		tags, response, err := client.Repositories.ListTags(context.TODO(), "kubernetes", "kubectl", listOptions)
		if response.NextPage == 0 {
			break
		}
		listOptions.Page = response.NextPage

		if err != nil {
			log.Fatal(err)
		}

		for _, tag := range tags {
			fmt.Println(tag.GetName())
		}
	}

	// Filters the references list and only keeps tags
	// var tags []string
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
		kctllist: []string{},
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
