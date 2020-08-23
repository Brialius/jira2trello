package trello

type Connector interface {
	Connect() error
	GetBoards() (map[string]*Board, error)
	GetLists() (map[string]*List, error)
	GetLabels() (map[string]*Label, error)
	GetUserJiraCards() ([]*Card, error)
	CreateCard(*Card) error
	MoveCardToList(string, string) error
	UpdateCardLabels(string, string) error
	SetBoard() error
	GetConfig() *Config
}
