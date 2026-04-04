package internal

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

var ErrSemverInvalid = errors.New("invalid semver")

type ErrSemverPrefixInvalid struct {
	GotPrefix string
}

func (e *ErrSemverPrefixInvalid) Error() string {
	return fmt.Sprintf(`expected prefix "v", got %q`, e.GotPrefix)
}

func (e *ErrSemverPrefixInvalid) Unwrap() error {
	return ErrSemverInvalid
}

type ErrSemverEmptyString struct{}

func (e *ErrSemverEmptyString) Error() string {
	return "semver cannot be an empty string"
}

func (e *ErrSemverEmptyString) Unwrap() error {
	return ErrSemverInvalid
}

// Takes a string which should be in the format v1.2.3 and returns a slice of strings. One string for each semver level: v1.2.3, v1.2, v1.
//
// If the input string is not a valid semver, then the returned slice will be empty.
func semverTags(src string) (tags []string, err error) {
	valid := semver.IsValid(src)
	if valid {
		tags = append(tags, src)
		min := semver.MajorMinor(src)
		if min != "" {
			tags = append(tags, min)
		}
		maj := semver.Major(src)
		if maj != "" {
			tags = append(tags, maj)
		}
	} else {
		if len(src) == 0 {
			err = &ErrSemverEmptyString{}
			return
		}
		if !strings.HasPrefix(src, "v") {
			err = &ErrSemverPrefixInvalid{GotPrefix: src[0:1]}
			return
		}
		err = fmt.Errorf("%w: %s", ErrSemverInvalid, src)
	}

	return
}
