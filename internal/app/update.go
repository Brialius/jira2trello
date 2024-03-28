package app

import (
	"context"
	"fmt"
	"github.com/creativeprojects/go-selfupdate"
	"os"
	"runtime"
)

const repositorySlug = "Brialius/jira2trello"

func DoSelfUpdate(currentVersion string) {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug(repositorySlug))
	if err != nil {
		fmt.Println("error occurred while detecting version:", err)

		return
	}

	if !found {
		fmt.Printf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)

		return
	}

	if latest.LessOrEqual(currentVersion) {
		fmt.Printf("Current version (%s) is the latest", currentVersion)

		return
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Println("could not locate executable path")

		return
	}

	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		fmt.Println("error occurred while updating binary:", err)

		return
	}

	fmt.Println("Successfully updated to version:", latest.Version())
}
