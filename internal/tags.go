package internal

import (
	"os"

	"golang.org/x/mod/semver"
)

// generate tags
// settings.tags are always added
// when auto tag is enabled, we will...
//   try to parse CI_COMMIT_TAG for version details
//     create tags for major, minor, patch levels
//     i.e. CI_COMMIT_TAG=v1.2.3 would have tags for
//       v1, v1.2, v1.2.3
//   whether parsing fails or succeeds, well set the CI_COMMIT_TAG as a tag

func (p *Plugin) generateTags() []string {
	tags := []string{}
	if p.Settings.AutoTag {
		ciCommitTag, found := os.LookupEnv("CI_COMMIT_TAG")
		if found {
			valid := semver.IsValid(ciCommitTag)
			if valid {
				maj := semver.Major(ciCommitTag)
				if maj != "" {
					tags = append(tags, maj)
				}
				min := semver.MajorMinor(ciCommitTag)
				if min != "" {
					tags = append(tags, min)
				}
			}
			tags = append(tags, ciCommitTag)
		}
	}

	return tags
}
