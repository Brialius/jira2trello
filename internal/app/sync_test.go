package app

import (
	"bytes"
	"encoding/json"
	"github.com/Brialius/jira2trello/internal/jira"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestSyncService_printJiraTasks(t *testing.T) {
	var jTasks = map[string]*jira.Task{}
	mustLoadJSONFile(t, "testdata/jira_tasks.json", &jTasks)

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
			wantOut: `Dependency               Story        JIRA1-984      Task name 984      18 May 20     0.0
ToDo                     Sub-task     JIRA1-990      Task name 990      18 May 20     0.0
ToDo                     Sub-task     JIRA1-991      Task name 991      18 May 20     0.5
In Dev / In Progress     Sub-task     JIRA1-987      Task name 987      28 May 20     1.0
In Dev / In Progress     Bug          JIRA1-223      Task name 223      15 Jun 20     1.0
In Dev / In Progress     Story        JIRA1-375      Task name 375      01 Jul 20     0.0
In Dev / In Progress     Sub-task     JIRA1-390      Task name 390      02 Jul 20     8.0
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
			s.printJiraTasks(out)
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
