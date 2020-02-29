/*
Copyright © 2019 Denis Belyatsky <denis.bel@gmail.com>

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
package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Brialius/jira2trello/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"sort"
	"strings"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Ask configuration settings and save them to file",
	Long:  `Ask configuration settings and save them to file`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := internal.MainConfig{}

		if err := viper.UnmarshalKey("users", &cfg.Users); err != nil {
			log.Printf("Can't get users from config: %s", err)
		}

		jiraQs := []*survey.Question{
			{
				Name: "url",
				Prompt: &survey.Input{
					Message: "What is jira server URL?",
					Default: viper.GetString("jira.url"),
				},
				Validate: survey.Required,
			},
			{
				Name: "user",
				Prompt: &survey.Input{
					Message: "What is jira username?",
					Default: viper.GetString("jira.user"),
				},
				Validate: survey.Required,
			},
			{
				Name: "password",
				Prompt: &survey.Password{
					Message: "What is jira password?",
				},
			},
		}
		_ = survey.Ask(jiraQs, &cfg.Jira)

		if cfg.Jira.Password == "" {
			cfg.Jira.Password = viper.GetString("jira.password")
		}

		viper.Set("jira", cfg.Jira)

		trelloQs := []*survey.Question{
			{
				Name: "apiKey",
				Prompt: &survey.Input{
					Message: "What is trello API key?",
					Default: viper.GetString("trello.apiKey"),
				},
				Validate: survey.Required,
			},
			{
				Name: "token",
				Prompt: &survey.Password{
					Message: "What is trello token?",
				},
			},
		}

		_ = survey.Ask(trelloQs, &cfg.Trello)

		if cfg.Trello.Token == "" {
			cfg.Trello.Token = viper.GetString("trello.token")
		}

		viper.Set("trello.apiKey", cfg.Trello.ApiKey)
		viper.Set("trello.token", cfg.Trello.Token)

		srv := internal.NewTrelloServer()
		err := srv.Connect()
		if err != nil {
			log.Fatalf("Can't connect to trello: %s", err)
		}

		boards, err := srv.GetBoards()
		if err != nil {
			log.Fatalf("Can't get trello boards: %s", err)
		}

		keys := make([]string, 0, len(boards))
		for k, v := range boards {
			keys = append(keys, k+" - "+v.Name)
		}
		sort.Strings(keys)

		var board string
		_ = survey.AskOne(&survey.Select{
			Message: "Please select trello board",
			Options: keys,
		}, &board)

		cfg.Trello.Board = board[:internal.TrelloIdLength]
		err = srv.SetBoard(board[:internal.TrelloIdLength])
		if err != nil {
			log.Fatalf("Can't set trello board: %s", err)
		}

		viper.Set("trello.board", &cfg.Trello.Board)

		lists, err := srv.GetLists()
		if err != nil {
			log.Fatalf("Can't get trello lists: %s", err)
		}

		keys = make([]string, 0, len(lists))
		for _, list := range lists {
			keys = append(keys, list.ID+" - "+list.Name)
		}
		sort.Strings(keys)

		var choice string

		_ = survey.AskOne(&survey.Select{
			Message: "Please select todo list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Lists.Todo = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select doing list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Lists.Doing = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select done list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Lists.Done = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select review list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Lists.Review = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select bucket list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Lists.Bucket = choice

		viper.Set("trello.lists", &cfg.Trello.Lists)

		labels, err := srv.GetLabels()
		if err != nil {
			log.Fatalf("Can't get trello labels: %s", err)
		}

		keys = make([]string, 0, len(labels))
		for _, label := range labels {
			keys = append(keys, label.ID+" - "+label.Name)
		}
		sort.Strings(keys)

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Jira label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Labels.Jira = choice[:internal.TrelloIdLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Blocked label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Labels.Blocked = choice[:internal.TrelloIdLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Task label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Labels.Task = choice[:internal.TrelloIdLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Bug label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Labels.Bug = choice[:internal.TrelloIdLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Story label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		cfg.Trello.Labels.Story = choice[:internal.TrelloIdLength]

		viper.Set("trello.labels", &cfg.Trello.Labels)

		members, err := srv.GetMembers()
		if err != nil {
			log.Fatalf("can't get members for %s board", board)
		}

		fmt.Printf("\nList of users from %s board:\n", board[internal.TrelloIdLength+2:])
		for _, member := range members {
			fmt.Printf("%s - %s(%s)\n", member.ID, member.FullName, member.Username)
		}

		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Can't write config: %s", err)
		}
		fmt.Println("Config updated")
	},
}

func removeKeyFromSlice(slice []string, k string) {
	for i, key := range slice {
		if strings.HasPrefix(key, k) {
			copy(slice[i:], slice[i+1:])
			slice[len(slice)-1] = ""
			slice = slice[:len(slice)-1]
			break
		}
	}
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
