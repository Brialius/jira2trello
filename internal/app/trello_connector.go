package app

import "github.com/Brialius/jira2trello/internal/trello"

//go:generate moq -out trello_connector_moq_test.go . TrelloConnector

type TrelloConnector interface {
	Connect() error
	GetBoards() (map[string]*trello.Board, error)
	GetLists() (map[string]*trello.List, error)
	GetLabels() (map[string]*trello.Label, error)
	GetUserJiraCards() ([]*trello.Card, error)
	CreateCard(*trello.Card) error
	MoveCardToList(string, string) error
	UpdateCardLabels(string, string) error
	SetBoard() error
	GetConfig() *trello.Config
}
