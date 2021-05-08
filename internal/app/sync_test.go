/*
Copyright © 2021 Denis Belyatsky <denis.bel@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package app

import (
	"bytes"
	"encoding/json"
	"github.com/Brialius/jira2trello/internal/jira"
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestSyncService_printJiraTasks(t *testing.T) {
	var jTasks = map[string]*jira.Task{}
	mustLoadJSONFile(t, "testdata/test_jira_tasks.json", &jTasks)

	type fields struct {
		jTasks map[string]*jira.Task
	}
	tests := []struct {
		name    string
		fields  fields
		wantOut string
	}{
		{
			name: "valid output",
			fields: fields{
				jTasks: jTasks,
			},
			wantOut: `In QA Review             Story        JIRA1-984      Task name 984      18 May 20     0.0
ToDo                     Sub-task     JIRA1-990      Task name 990      18 May 20     0.0
ToDo                     Sub-task     JIRA1-991      Task name 991      18 May 20     0.5
In Dev / In Progress     Sub-task     JIRA1-987      Task name 987      28 May 20     1.0
In Dev / In Progress     Bug          JIRA1-223      Task name 223      15 Jun 20     1.0
In Dev / In Progress     Story        JIRA1-375      Task name 375      01 Jul 20     0.0
In Dev / In Progress     Sub-task     JIRA1-391      Task name 391      02 Jul 20     10.0
In Dev / In Progress     Sub-task     JIRA1-392      Task name 392      02 Jul 20     6.0
In Dev / In Progress     Sub-task     JIRA1-431      Task name 431      08 Jul 20     12.0
In Dev / In Progress     Sub-task     JIRA1-433      Task name 433      08 Jul 20     19.0
In Dev / In Progress     Sub-task     JIRA1-434      Task name 434      08 Jul 20     12.0
ToDo                     Sub-task     JIRA1-1110     Task name 1110     07 Aug 20     11.0
In Dev / In Progress     Sub-task     JIRA1-1130     Task name 1130     08 Aug 20     13.0
In Dev / In Progress     Sub-task     JIRA1-1131     Task name 1131     08 Aug 20     14.0
In Dev / In Progress     Sub-task     JIRA1-1133     Task name 1133     10 Aug 20     2.0
In Dev / In Progress     Bug          JIRA1-1194     Task name 1194     12 Aug 20     1.0
In Dev / In Progress     Task         JIRA1-1195     Task name 1195     12 Aug 20     1.0
In Dev / In Progress     Story        JIRA1-1288     Task name 1288     19 Aug 20     0.0
In Dev / In Progress     Task         JIRA1-1304     Task name 1304     20 Aug 20     2.0
ToDo                     Sub-task     JIRA1-1324     Task name 1324     21 Aug 20     0.0
`,
		},
		{
			name: "empty output",
			fields: fields{
				jTasks: map[string]*jira.Task{},
			},
			wantOut: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				jTasks: tt.fields.jTasks,
			}
			out := &bytes.Buffer{}
			printJiraTasks(colorable.NewNonColorable(out), s.jTasks)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printJiraTasks() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func mustLoadJSONFile(t *testing.T, file string, variable interface{}) []byte {
	testFileContent, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	err = json.Unmarshal(testFileContent, &variable)
	require.NoError(t, err)

	testFileContent, err = json.Marshal(variable)
	require.NoError(t, err)

	return testFileContent
}

func TestSyncService_Sync(t *testing.T) {
	var jTasks = map[string]*jira.Task{}
	mustLoadJSONFile(t, "testdata/test_jira_tasks.json", &jTasks)

	tCards := make([]*trello.Card, 0)
	mustLoadJSONFile(t, "testdata/test_trello_cards.json", &tCards)

	jCli := GetJiraMockedCli(jTasks)
	tCli := GetTrelloMockedCli(tCards)

	type fields struct {
		jCli   JiraConnector
		tCli   TrelloConnector
		jTasks map[string]*jira.Task
		tCards map[string]*trello.Card
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "valid",
			fields: fields{
				jCli:   jCli,
				tCli:   tCli,
				jTasks: nil,
				tCards: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				jCli:   tt.fields.jCli,
				tCli:   tt.fields.tCli,
				jTasks: tt.fields.jTasks,
				tCards: tt.fields.tCards,
			}
			s.Sync()

			require.Equal(t, []struct{ S1, S2 string }{
				{"098098098098098098098011", "121212121212121212121fa4,12121212121212121212a0c8"}},
				tCli.UpdateCardLabelsCalls())

			require.Equal(t, []struct{ S1, S2 string }{
				{"098098098098098098098011", "12345678909876543219d1cc"},
				{"098098098098098098098018", "12345678909876543219d1cf"},
			},
				tCli.MoveCardToListCalls())

			require.Equal(t, []struct{ Card *trello.Card }{{Card: &trello.Card{
				Name:      "JIRA1-1194 | Task name 1194",
				ListID:    "12345678909876543219d1cb",
				Desc:      "\nJira link: https://jira-site/browse/JIRA1-1194\nType: Bug",
				IDLabels:  &[]string{"121212121212121212121fa4", "12121212121212121212de33"},
				IDMembers: "111111111111111111111111",
			}}}, tCli.CreateCardCalls())
		})
	}
}

func GetJiraMockedCli(jTasks map[string]*jira.Task) *JiraConnectorMock {
	return &JiraConnectorMock{
		ConnectFunc: func() error {
			return nil
		},
		GetUserTasksFunc: func(jql string) (map[string]*jira.Task, error) {
			return jTasks, nil
		},
	}
}

func GetTrelloMockedCli(tCards []*trello.Card) *TrelloConnectorMock {
	return &TrelloConnectorMock{
		ConnectFunc: func() error {
			return nil
		},
		CreateCardFunc: func(in1 *trello.Card) error {
			return nil
		},
		GetBoardsFunc: func() (map[string]*trello.Board, error) {
			return map[string]*trello.Board{
				"Board1": {
					"https://trello.com/b/0/board1",
					"Board",
					"000000000000000000000000",
				},
				"Board2": {
					"https://trello.com/b/1/board2",
					"Board",
					"111111111111111111111111",
				},
				"Board3": {
					"https://trello.com/b/2/board3",
					"Board",
					"222222222222222222222222",
				},
				"Board4": {
					"https://trello.com/b/3/board4",
					"Board",
					"33333333333333333333333",
				},
			}, nil
		},
		GetConfigFunc: func() *trello.Config {
			return &trello.Config{
				UserID: "111111111111111111111111",
				Lists: &trello.Lists{
					Todo:   "12345678909876543219d1c9",
					Doing:  "12345678909876543219d1cb",
					Done:   "12345678909876543219d1cf",
					Review: "12345678909876543219d1cc",
					Bucket: "12345678909876543219d1d0",
				},
				Labels: &trello.Labels{
					Jira:    "121212121212121212121fa4",
					Blocked: "12121212121212121212d298",
					Task:    "121212121212121212121795",
					Bug:     "12121212121212121212de33",
					Story:   "12121212121212121212a0c8",
				},
				Debug: false,
			}
		},
		GetLabelsFunc: func() (map[string]*trello.Label, error) {
			return map[string]*trello.Label{
				"Jira":    {"Jira", "121212121212121212121fa4"},
				"Blocked": {"Blocked", "12121212121212121212d298"},
				"Task":    {"Task", "121212121212121212121795"},
				"Bug":     {"Bug", "12121212121212121212de33"},
				"Story":   {"Story", "12121212121212121212a0c8"},
			}, nil
		},
		GetListsFunc: func() (map[string]*trello.List, error) {
			return map[string]*trello.List{
				"Todo":   {"Todo", "12345678909876543219d1c9"},
				"Doing":  {"Doing", "12345678909876543219d1cb"},
				"Done":   {"Done", "12345678909876543219d1cf"},
				"Review": {"Review", "12345678909876543219d1cc"},
				"Bucket": {"Bucket", "12345678909876543219d1d0"},
			}, nil
		},
		GetUserJiraCardsFunc: func() ([]*trello.Card, error) {
			return tCards, nil
		},
		MoveCardToListFunc: func(in1 string, in2 string) error {
			return nil
		},
		SetBoardFunc: func() error {
			return nil
		},
		UpdateCardLabelsFunc: func(in1 string, in2 string) error {
			return nil
		},
	}
}
