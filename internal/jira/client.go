/*
Copyright © 2021 Denis Belyatsky <denis.bel@gmail.com>

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
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"time"
)

type Client struct {
	*Config
	cli *jira.Client
}

func NewClient(cfg *Config) *Client {
	return &Client{
		Config: cfg,
	}
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

func (j *Client) GetUserTasks(jql string) (map[string]*Task, error) {
	res := map[string]*Task{}
	issues, _, err := j.cli.Issue.Search("assignee = currentUser() AND "+jql, nil)

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

	j.writeToJSONFile(res, "jira_tasks.json")

	return res, nil
}

func (j *Client) writeToJSONFile(value interface{}, fileName string) {
	if j.Debug {
		b, _ := json.MarshalIndent(value, "", "  ")
		err := ioutil.WriteFile(fileName, b, 0600)

		if err != nil {
			fmt.Printf("can't write debug file: %s", err)
		}
	}
}
