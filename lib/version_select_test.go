package lib_test

import (
	"kctlswitch/lib"
	"testing"

	"github.com/Masterminds/semver/v3"
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

var invalidConstraints = []struct {
	name       string
	constraint string
	want       error
}{
	{"invalid constraint", "vers1.2.3.4.5-1", semver.ErrInvalidCharacters},
	{"syntax valid semantic invalid constraint", "v0.1.1", semver.ErrSegmentStartsZero},
	{"empty constraint", "", semver.ErrEmptyString},
}

func TestVersionListValidConstraints(t *testing.T) {
	for _, ct := range validConstraints {
		t.Run(ct.name, func(t *testing.T) {
			got, err := lib.KctlVersionList(ct.constraint, slog)
			assert.NoError(t, err)
			assert.Equal(t, ct.want, got)
		})
	}
}

func TestVersionListInvalidConstraints(t *testing.T) {
	for _, ict := range invalidConstraints {
		t.Run(ict.name, func(t *testing.T) {
			v, err := lib.KctlVersionList(ict.constraint, slog)
			assert.Empty(t, v)
			assert.Error(t, ict.want, err)
		})
	}
}
