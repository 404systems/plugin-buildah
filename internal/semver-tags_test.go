package internal

import (
	"errors"
	"reflect"
	"testing"
)

func TestGenerateTags(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   []string
		shouldFail bool
	}{
		{
			name:  "valid semver",
			input: "v1.2.3",
			expected: []string{
				"v1.2.3",
				"v1.2",
				"v1",
			},
			shouldFail: false,
		},
		{
			name:       "invalid semver",
			input:      "1.2.3",
			expected:   []string(nil),
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := semverTags(tt.input)
			if err != nil {
				if errors.Is(err, ErrSemverInvalid) {
					if !tt.shouldFail {
						t.Fatalf("failed to parse semver: %v", err)
					}
				}
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("got %q, wanted %q", got, tt.expected)
			}
			t.Logf("got %q", got)
		})
	}
}
