package trello

import "fmt"

const IdLength = 24

type Lists struct {
	Todo   string
	Doing  string
	Done   string
	Review string
	Bucket string
}

type Labels struct {
	Jira    string
	Blocked string
	Task    string
	Bug     string
	Story   string
}

type Card struct {
	ID        string
	Name      string
	ListID    string
	List      string
	Labels    string
	Key       string
	Desc      string
	IDLabels  *[]string
	IDMembers string
}

type Board struct {
	URL  string
	Name string
	ID   string
}

type Label struct {
	Name string
	ID   string
}

type List struct {
	Name string
	ID   string
}

type Member struct {
	Name     string
	FullName string
	ID       string
}

func (c *Card) String() string {
	return fmt.Sprintf("%s | %s(%s): %s - %s", c.Key, c.List, c.ListID, c.Name, *c.IDLabels)
}
