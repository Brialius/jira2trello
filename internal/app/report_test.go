package app

import (
	"bytes"
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/mattn/go-colorable"
	"testing"
)

func Test_printReport(t *testing.T) {
	var tCards = make([]*trello.Card, 0)
	mustLoadJSONFile(t, "testdata/test_trello_cards.json", &tCards)

	tCli := GetTrelloMockedCli(tCards)

	type args struct {
		tCli   TrelloConnector
		tCards []*trello.Card
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "valid",
			args: args{
				tCli:   tCli,
				tCards: tCards,
			},
			wantOut: `
----------------------------------

JIRA1-1289 | Test task 1289 - Done
https://jira-site/browse/JIRA1-1289

JIRA1-1290 | Test task 1290 - Done
https://jira-site/browse/JIRA1-1290

JIRA1-984 | Test task 984 - In progress
https://jira-site/browse/JIRA1-984

JIRA1-1133 | Test task 1133 - In progress
https://jira-site/browse/JIRA1-1133

JIRA1-1288 | Test task 1288 - In progress
https://jira-site/browse/JIRA1-1288

JIRA1-1130 | Test task 1130 - In progress
https://jira-site/browse/JIRA1-1130

JIRA1-1131 | Test task 1131 - In progress
https://jira-site/browse/JIRA1-1131

JIRA1-987 | Test task 987 - In progress
https://jira-site/browse/JIRA1-987

JIRA1-1195 | Test task 1195 - In progress
https://jira-site/browse/JIRA1-1195

JIRA1-223 | Test task 223 - In progress
https://jira-site/browse/JIRA1-223

JIRA1-375 | Test task 375 - In progress
https://jira-site/browse/JIRA1-375

JIRA1-434 | Test task 434 - In progress
https://jira-site/browse/JIRA1-434

JIRA1-433 | Test task 433 - In progress
https://jira-site/browse/JIRA1-433

JIRA1-431 | Test task 431 - In progress
https://jira-site/browse/JIRA1-431

JIRA1-392 | Test task 392 - In progress
https://jira-site/browse/JIRA1-392

JIRA1-390 | Test task 390 - In progress
https://jira-site/browse/JIRA1-390

JIRA1-391 | Test task 391 - In progress
https://jira-site/browse/JIRA1-391

JIRA1-1304 | Test task 1304 - In progress
https://jira-site/browse/JIRA1-1304

JIRA1-1324 | Test task 1324 - In review
https://jira-site/browse/JIRA1-1324

----------------------------------
In progress: 16
In review: 1
Done: 2
`,
		},
		{
			name: "empty",
			args: args{
				tCli:   tCli,
				tCards: []*trello.Card{},
			},
			wantOut: `
----------------------------------

----------------------------------
In progress: 0
In review: 0
Done: 0
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printReport(colorable.NewNonColorable(out), tt.args.tCli, tt.args.tCards, "https://jira-site")
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printReport() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
