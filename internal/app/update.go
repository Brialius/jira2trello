package app

import (
	"fmt"
	"github.com/Brialius/jira2trello/internal"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"strings"
)

func DoSelfUpdate() {
	v := semver.MustParse(strings.TrimPrefix(internal.Version, "v"))

	slug := "Brialius/jira2trello"
	latest, found, err := selfupdate.DetectLatest(slug)

	if err != nil {
		fmt.Println("Failed to check updates:", err)

		return
	}

	if latest.Version.Equals(v) {
		// latest version is the same as current version. It means current binary is up to date.
		fmt.Println("Current binary is the latest version:", internal.Version)

		return
	}

	if found && latest.Version.GT(v) {
		fmt.Printf("New version found: %s\n", latest.Version)
		fmt.Println("Updating...")

		_, err := selfupdate.UpdateSelf(v, slug)

		if err != nil {
			fmt.Println("Binary update failed:", err)

			return
		}

		fmt.Println("Successfully updated to version", latest.Version)
		fmt.Println("Release notes:")
		fmt.Println(latest.ReleaseNotes)
	}
}
