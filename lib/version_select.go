package lib

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/semver"
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
