package internal

import (
	"reflect"
	"testing"
)

func TestBuildTimeArgValueUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected buildTimeArgValue
	}{
		{
			name:     "quoted string",
			input:    []byte(`"hello world"`),
			expected: buildTimeArgValue("hello world"),
		},
		{
			name:     "unquoted number",
			input:    []byte(`42`),
			expected: buildTimeArgValue("42"),
		},
		{
			name:     "unquoted bool",
			input:    []byte(`true`),
			expected: buildTimeArgValue("true"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got buildTimeArgValue
			if err := got.UnmarshalJSON(tt.input); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Fatalf("got %s, wanted %s", got, tt.expected)
			}
		})
	}
}

func TestBuildTimeArgParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "version data with commit as string",
			input: `{"VERSION":"devel","COMMIT":"abc1234"}`,
			expected: map[string]string{
				"VERSION": "devel",
				"COMMIT":  "abc1234",
			},
		},
		{
			name:  "version data with commit as number",
			input: `{"VERSION":"devel","COMMIT":1234567}`,
			expected: map[string]string{
				"VERSION": "devel",
				"COMMIT":  "1234567",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBuildTimeArgs(tt.input)
			if err != nil {
				t.Fatalf("failed to parse build time args: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("got %q, wanted %q", got, tt.expected)
			}
		})
	}
}
