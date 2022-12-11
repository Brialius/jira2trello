package app

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"strings"
)

func DoSelfUpdate(currentVersion string) {
	ver := semver.MustParse(strings.TrimPrefix(currentVersion, "v"))

	slug := "Brialius/jira2trello"
	latest, found, err := selfupdate.DetectLatest(slug)

	if err != nil {
		fmt.Println("Failed to check updates:", err)

		return
	}

	if latest.Version.Equals(ver) {
		// latest version is the same as current version. It means current binary is up-to-date.
		fmt.Println("Current binary is the latest version:", currentVersion)

		return
	}

	if found && latest.Version.GT(ver) {
		fmt.Printf("New version found: %s\n", latest.Version)
		fmt.Println("Updating...")

		_, err := selfupdate.UpdateSelf(ver, slug)

		if err != nil {
			fmt.Println("Binary update failed:", err)

			return
		}

		fmt.Println("Successfully updated to version", latest.Version)
		fmt.Println("Release notes:")
		fmt.Println(latest.ReleaseNotes)
	}
}
