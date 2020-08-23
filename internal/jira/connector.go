package jira

type Connector interface {
	Connect() error
	GetUserTasks() (map[string]*Task, error)
}
