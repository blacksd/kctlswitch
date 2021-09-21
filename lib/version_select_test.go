package lib_test

import (
	"kctlswitch/lib"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConstraints = []struct {
	name       string
	constraint string
	want       []string
}{
	{"single constraint", "1.18.4", []string{"1.18.4"}},
	{"range constraint", "1.17.5 - 1.17.9", []string{"1.17.5", "1.17.6", "1.17.7", "1.17.8", "1.17.9"}},
	{"constraint that yields zero results", "1.10.143", []string{}},
}

/* var invalidConstraints = []struct {
	name       string
	constraint string
	err        error
}{
	{"invalid constraint", "vers0.2.3.4.5", errors.Error()},
	{"empty constraint", "", errors.Error()},
} */

func TestVersionListValidConstraints(t *testing.T) {
	for _, ct := range validConstraints {
		t.Run(ct.name, func(t *testing.T) {
			got, err := lib.KctlVersionList(ct.constraint, slog)
			assert.NoError(t, err)
			assert.Equal(t, ct.want, got)
		})
	}
}

/* t.Run("InvalidConstraint", func(t *testing.T) {
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
}) */
