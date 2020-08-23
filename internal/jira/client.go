/*
Copyright © 2019 Denis Belyatsky <denis.bel@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package jira

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
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

type Client struct {
	Config
	cli *jira.Client
}

func NewServer(cfg Config) *Client {
	return &Client{
		Config: Config{
			User:     cfg.User,
			Password: cfg.Password,
			URL:      cfg.URL,
		},
	}
}

func (j Task) String() string {
	return fmt.Sprintf("%s | %s | %s | %s, %s, (%0.1f)",
		j.Status, j.Type, j.Key, j.Summary, j.Created.Format(time.RFC822), j.TimeSpent.Hours())
}

func (j Task) TabString() string {
	return fmt.Sprintf("%s \t%s \t%s \t%.70s \t%.9s \t%0.1f",
		j.Status, j.Type, j.Key, j.Summary, j.Created.Format(time.RFC822), j.TimeSpent.Hours())
}

func (j *Client) Connect() error {
	tp := jira.BasicAuthTransport{
		Username: j.User,
		Password: j.Password,
	}

	client, err := jira.NewClient(tp.Client(), j.URL)
	if err != nil {
		return err
	}

	j.cli = client

	return nil
}

func (j *Client) GetUserTasks() (map[string]*Task, error) {
	res := map[string]*Task{}
	issues, _, err := j.cli.Issue.Search("assignee = '"+j.User+"' AND status not in "+
		"(done, closed, close, resolved) ORDER BY priority DESC, updated DESC", nil)

	if err != nil {
		return nil, err
	}

	for _, issue := range issues {
		res[issue.Key] = &Task{
			Created:   time.Time(issue.Fields.Created),
			Updated:   time.Time(issue.Fields.Updated),
			TimeSpent: time.Duration(issue.Fields.TimeSpent) * time.Second,
			Summary:   issue.Fields.Summary,
			Link:      j.URL + "/browse/" + issue.Key,
			Self:      issue.Self,
			Key:       issue.Key,
			Status:    issue.Fields.Status.Name,
			Desc:      issue.Fields.Description,
			Type:      issue.Fields.Type.Name,
		}
		if parent := issue.Fields.Parent; parent != nil {
			res[issue.Key].ParentID = parent.ID
			res[issue.Key].ParentKey = parent.Key
			res[issue.Key].ParentLink = j.URL + "/browse/" + parent.Key
		}
	}

	return res, nil
}