package app

import (
	"bytes"
	"github.com/Brialius/jira2trello/internal/trello"
	"testing"
)

func Test_printReport(t *testing.T) {
	var tCards = make([]*trello.Card, 0)
	mustLoadJSONFile(t, "testdata/trello_cards.json", &tCards)

	tCli := trello.NewTestClient(&trello.Config{
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
	})

	type args struct {
		tCli   trello.Connector
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
JIRA1-1304 | Test task 1304 - In progress
https://jira.inbcu.com/browse/JIRA1-1304
JIRA1-1194 | Test task 1194 - In progress
https://jira.inbcu.com/browse/JIRA1-1194
JIRA1-1195 | Test task 1195 - In progress
https://jira.inbcu.com/browse/JIRA1-1195
JIRA1-1133 | Test task 1133 - In progress
https://jira.inbcu.com/browse/JIRA1-1133
JIRA1-1288 | Test task 1288 - In progress
https://jira.inbcu.com/browse/JIRA1-1288
JIRA1-1130 | Test task 1130 - In progress
https://jira.inbcu.com/browse/JIRA1-1130
JIRA1-1131 | Test task 1131 - In progress
https://jira.inbcu.com/browse/JIRA1-1131
JIRA1-987 | Test task 987 - In progress
https://jira.inbcu.com/browse/JIRA1-987
JIRA1-984 | Test task 984 - In progress
https://jira.inbcu.com/browse/JIRA1-984
JIRA1-223 | Test task 223 - In progress
https://jira.inbcu.com/browse/JIRA1-223
JIRA1-375 | Test task 375 - In progress
https://jira.inbcu.com/browse/JIRA1-375
JIRA1-434 | Test task 434 - In progress
https://jira.inbcu.com/browse/JIRA1-434
JIRA1-433 | Test task 433 - In progress
https://jira.inbcu.com/browse/JIRA1-433
JIRA1-431 | Test task 431 - In progress
https://jira.inbcu.com/browse/JIRA1-431
JIRA1-392 | Test task 392 - In progress
https://jira.inbcu.com/browse/JIRA1-392
JIRA1-390 | Test task 390 - In progress
https://jira.inbcu.com/browse/JIRA1-390
JIRA1-391 | Test task 391 - In progress
https://jira.inbcu.com/browse/JIRA1-391
JIRA1-1290 | Test task 1290 - Done
https://jira.inbcu.com/browse/JIRA1-1290
JIRA1-1289 | Test task 1289 - Done
https://jira.inbcu.com/browse/JIRA1-1289

----------------------------------
In progress: 17
In review: 0
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
			printReport(out, tt.args.tCli, tt.args.tCards)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printReport() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
