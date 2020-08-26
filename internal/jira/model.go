package jira

import (
	"fmt"
	"github.com/Brialius/jira2trello/internal"
	"time"
)

type Task struct {
	Created    time.Time
	Updated    time.Time
	TimeSpent  time.Duration
	Summary    string
	Link       string
	Self       string
	Key        string
	Status     string
	Desc       string
	ParentID   string
	ParentKey  string
	ParentLink string
	Type       string
}

func (j Task) String() string {
	return fmt.Sprintf("%s | %s | %s | %s, %s, (%0.1f)",
		j.Status, j.Type, j.Key, j.Summary, j.Created.Format(time.RFC822), j.TimeSpent.Hours())
}

func (j Task) TabString() string {
	jStatus := j.Status
	jType := j.Type

	switch jStatus {
	case "Dependency", "Blocked":
		jStatus = internal.Red + jStatus + internal.ColorOff
	case "ToDo":
		jStatus = internal.Blue + jStatus + internal.ColorOff
	case "In QA Review":
		jStatus = internal.Cyan + jStatus + internal.ColorOff
	default:
		jStatus = internal.Yellow + jStatus + internal.ColorOff
	}

	switch jType {
	case "Story", "User Story":
		jType = internal.Green + jType + internal.ColorOff
	case "Bug":
		jType = internal.Red + jType + internal.ColorOff
	default:
		jType = internal.Blue + jType + internal.ColorOff
	}

	return fmt.Sprintf("%s \t%s \t%s \t%.70s \t%.9s \t%0.1f",
		jStatus, jType, j.Key, j.Summary, j.Created.Format(time.RFC822), j.TimeSpent.Hours())
}
