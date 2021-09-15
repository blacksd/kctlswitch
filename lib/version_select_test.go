package lib_test

import (
	"kctlswitch/lib"
	"strings"
	"testing"
)

const (
	goodConstraint string = ">= 1.18"
	badConstraint  string = "vers0.2.3.4.5"
)

func TestFetchGitTags(t *testing.T) {
	t.Run("SingleValidConstraint", func(t *testing.T) {
		sample, err := lib.FetchGitTags(goodConstraint)
		if err != nil {
			t.Fatal("Couldn't obtain data from a valid constraint.")
		}
		if len(sample) == 0 {
			t.Fatal("Got zero results for a valid constraint.")
		}
	})

	t.Run("InvalidConstraint", func(t *testing.T) {
		_, err := lib.FetchGitTags(badConstraint)
		if !strings.Contains(err.Error(), "improper constraint") {
			t.Fatal("An invalid constraint did not failed.")
		}

	})
}

/*
func TestValidateTag(t *testing.T) {

	t.Run("ValidConstraint", func(t *testing.T) {
		c, _ := semver.NewConstraint(goodConstraint)
		if err := validateTag(*c); err != nil {
			t.Fatal("Can't validate a good constraint.")
		}
	})
	t.Run("EmptyConstraint", func(t *testing.T) {
		c, _ := semver.NewConstraint("")
		if err := validateTag(*c); err != nil {
			t.Fatal("Can't validate an empty constraint.")
		}
	})
}
*/
