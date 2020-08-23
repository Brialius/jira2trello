package trello

import (
	"github.com/adlio/trello"
)

type Connector interface {
	Connect() error
	GetBoards() (map[string]*Board, error)
	GetLists() ([]*List, error)
	GetLabels() ([]*Label, error)
	GetBoardByID(string) (*Board, error)
	GetMembers() ([]*trello.Member, error)
	GetCards() ([]*Card, error)
	CreateCard(*Card) error
	MoveCardToList(string, string) error
	UpdateCardLabels(string, string) error
	SetBoard(string) error
	GetConfig() *Config
}