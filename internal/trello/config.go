package trello

type Config struct {
	ApiKey string
	Token  string
	Board  string
	Lists  Lists
	Labels Labels
}
