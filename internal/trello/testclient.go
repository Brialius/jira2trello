package trello

type TestClient struct {
	*Config
}

func NewTestClient(cfg *Config) *TestClient {
	return &TestClient{
		Config: cfg,
	}
}

func (t *TestClient) Connect() error {
	return nil
}

func (t *TestClient) GetBoards() (map[string]*Board, error) {
	return nil, nil
}

func (t *TestClient) GetLists() (map[string]*List, error) {
	return nil, nil
}

func (t *TestClient) GetLabels() (map[string]*Label, error) {
	return nil, nil
}

func (t *TestClient) GetUserJiraCards() ([]*Card, error) {
	return nil, nil
}

func (t *TestClient) CreateCard(card *Card) error {
	return nil
}

func (t *TestClient) MoveCardToList(s string, s2 string) error {
	return nil
}

func (t *TestClient) UpdateCardLabels(s string, s2 string) error {
	return nil
}

func (t *TestClient) SetBoard() error {
	return nil
}

func (t *TestClient) GetConfig() *Config {
	return t.Config
}
