package jira

import (
	"fmt"
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
	return fmt.Sprintf("%s \t%s \t%s \t%.70s \t%.9s \t%0.1f",
		j.Status, j.Type, j.Key, j.Summary, j.Created.Format(time.RFC822), j.TimeSpent.Hours())
}
