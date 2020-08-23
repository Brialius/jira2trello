package trello

type Config struct {
	APIKey string
	Token  string
	Board  string
	UserID string
	Lists  Lists
	Labels Labels
}
