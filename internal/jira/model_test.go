package jira

import (
	"testing"
	"time"
)

func TestTask_String(t *testing.T) {
	type fields struct {
		Created    time.Time
		Updated    time.Time
		DueDate    time.Time
		TimeSpent  time.Duration
		Summary    string
		Link       string
		Self       string
		Key        string
		Status     string
		Desc       string
		ParentID   string
		ParentKey  string
		ParentLink string
		Type       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid",
			fields: fields{
				Created:   time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC),
				Updated:   time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC),
				DueDate:   time.Date(2020, 1, 2, 1, 1, 1, 1, time.UTC),
				TimeSpent: 36000000000000,
				Summary:   "Test task 132",
				Link:      "https://jira-site/browse/JIRA1-132",
				Key:       "JIRA1-132",
				Status:    "In review",
				Type:      "Task",
			},
			want: "In review | Task | JIRA1-132 | Test task 132, 01 Jan 20 01:01 UTC, 02 Jan 20 01:01 UTC, (10.0)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Task{
				Created:    tt.fields.Created,
				Updated:    tt.fields.Updated,
				DueDate:    tt.fields.DueDate,
				TimeSpent:  tt.fields.TimeSpent,
				Summary:    tt.fields.Summary,
				Link:       tt.fields.Link,
				Self:       tt.fields.Self,
				Key:        tt.fields.Key,
				Status:     tt.fields.Status,
				Desc:       tt.fields.Desc,
				ParentID:   tt.fields.ParentID,
				ParentKey:  tt.fields.ParentKey,
				ParentLink: tt.fields.ParentLink,
				Type:       tt.fields.Type,
			}
			if got := j.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
