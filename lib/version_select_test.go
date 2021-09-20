package lib_test

import (
	"kctlswitch/lib"
	"strings"
	"testing"
)

const (
	goodConstraint  string = ">= 1.18"
	badConstraint   string = "vers0.2.3.4.5"
	emptyConstraint string = ""
)

func TestVersionList(t *testing.T) {
	t.Run("SingleValidConstraint", func(t *testing.T) {
		sample, err := lib.KctlVersionList(goodConstraint, slog)
		if err != nil {
			t.Fatal("Couldn't obtain data from a valid constraint.")
		}
		if len(sample) == 0 {
			t.Fatal("Got zero results for a valid constraint.")
		}
	})

	t.Run("InvalidConstraint", func(t *testing.T) {
		_, err := lib.KctlVersionList(badConstraint, slog)
		if !strings.Contains(err.Error(), "improper constraint") {
			t.Fatal("An invalid constraint did not failed.")
		}
	})

	t.Run("EmptyConstraint", func(t *testing.T) {
		_, err := lib.KctlVersionList(emptyConstraint, slog)
		if !strings.Contains(err.Error(), "improper constraint") {
			t.Fatal("The empty constraint was not managed.")
		}
	})
}
