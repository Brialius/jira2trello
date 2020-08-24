package app

import "github.com/Brialius/jira2trello/internal/jira"

//go:generate moq -out jira_connector_moq_test.go . JiraConnector

type JiraConnector interface {
	Connect() error
	GetUserTasks() (map[string]*jira.Task, error)
}
