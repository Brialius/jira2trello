package trello

type Config struct {
	APIKey string
	Token  string
	Board  string
	Lists  Lists
	Labels Labels
}
