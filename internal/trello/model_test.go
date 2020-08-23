package trello

import "testing"

func TestCard_IsInAnyOfLists(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		ListID    string
		List      string
		Labels    string
		Key       string
		Desc      string
		IDLabels  *[]string
		IDMembers string
	}
	type args struct {
		lists []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				ListID: "123456789012345678901234",
			},
			args: args{
				lists: []string{
					"123456789012345678901234",
					"123456789012345678901237",
				},
			},
			want: true,
		},
		{
			name: "false",
			fields: fields{
				ListID: "123456789012345678901234",
			},
			args: args{
				lists: []string{
					"123456789012345678901235",
					"123456789012345678901236",
					"123456789012345678901237",
				},
			},
			want: false,
		},
		{
			name: "long",
			fields: fields{
				ListID: "123456789012345678901235",
			},
			args: args{
				lists: []string{
					"123456789012345678901234 - Issues and Requests",
					"123456789012345678901235 - Done",
					"123456789012345678901236 - Doing",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Card{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				ListID:    tt.fields.ListID,
				List:      tt.fields.List,
				Labels:    tt.fields.Labels,
				Key:       tt.fields.Key,
				Desc:      tt.fields.Desc,
				IDLabels:  tt.fields.IDLabels,
				IDMembers: tt.fields.IDMembers,
			}
			if got := c.IsInAnyOfLists(tt.args.lists); got != tt.want {
				t.Errorf("IsInAnyOfLists() = %v, want %v", got, tt.want)
			}
		})
	}
}