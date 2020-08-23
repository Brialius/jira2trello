package jira

type TestClient struct {
	Config
	jTasks      map[string]*Task
	returnError error
}

func (t *TestClient) Connect() error {
	return nil
}

func (t *TestClient) GetUserTasks() (map[string]*Task, error) {
	return t.jTasks, t.returnError
}
