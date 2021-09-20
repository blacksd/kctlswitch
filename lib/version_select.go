package lib

import (
	"errors"

	"github.com/Masterminds/semver/v3"
	"go.uber.org/zap"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func KctlVersionList(constraint string, log *zap.SugaredLogger) ([]string, error) {
	var tags []string
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		log.Errorf("The constraint \"%s\" is not valid.", constraint)
		return nil, err
	}

	// Create the remote with repository URL
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/kubernetes/kubernetes.git"},
	})

	log.Debug("Fetching tags")

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Filters the references list and only keeps tags
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tag := ref.Name().Short()
			if err := validation.Validate(tag, validation.By(validateTag(*c))); err == nil {
				tags = append(tags, tag)
			}
		}
	}

	if len(tags) == 0 {
		log.Error("No tags are satisfying the constraint!")
	}

	log.Debugf("Tags found: %v", tags)

	return tags, nil
}

func validateTag(constraint semver.Constraints) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		ver, err := semver.NewVersion(s)
		if err != nil {
			return errors.New("The tag '" + s + "' is not semver compliant.")
		}
		if !constraint.Check(ver) {
			return errors.New("The tag '" + s + "' is not within the constraints.")
		}
		return nil
	}
}
