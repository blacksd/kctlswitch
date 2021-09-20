package lib_test

import (
	"kctlswitch/lib"
	"strings"
	"testing"
)

var validConstraints = []struct {
	name       string
	constraint string
	want       []string
}{
	{"single constraint", "1.18", TODO},
	{"range constraint", "1.17 - 1.18", TODO},
}

var invalidConstraints = []struct {
	name       string
	constraint string
	err        error
}{
	{"bad constraint", "vers0.2.3.4.5", errors.Error()},
	{"empty constraints", "", errors.Error()},
}

func TestVersionListValidConstraints(t *testing.T) {
	for _, ct := range validConstraints {
		t.Run(ct.name, func(t *testing.T) {
			got, err := lib.KctlVersionList(ct.constraint, slog)
			if err != nil {
				t.Fatal("Couldn't obtain data.")
			}
			if len(got) == 0 {
				t.Fatal("Got zero results.")
			}
		})
	}

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
