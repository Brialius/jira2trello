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
	"github.com/Brialius/jira2trello/internal/app"
	"github.com/Brialius/jira2trello/internal/jira"
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"sort"
	"strings"
)

// configureCmd represents the configure command.
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Ask configuration settings and save them to file",
	Long:  `Ask configuration settings and save them to file`,
	Run: func(cmd *cobra.Command, args []string) {
		config := app.Config{}

		if err := viper.UnmarshalKey("users", &config.Users); err != nil {
			log.Printf("Can't get users from config: %s", err)
		}

		jiraConfig := jira.Config{}

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
		_ = survey.Ask(jiraQs, &jiraConfig)

		if jiraConfig.Password == "" {
			jiraConfig.Password = viper.GetString("jira.password")
		}

		viper.Set("jira", jiraConfig)

		trelloConfig := trello.Config{}

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

		_ = survey.Ask(trelloQs, &trelloConfig)

		if trelloConfig.Token == "" {
			trelloConfig.Token = viper.GetString("trello.token")
		}

		viper.Set("trello.apiKey", trelloConfig.APIKey)
		viper.Set("trello.token", trelloConfig.Token)

		tSrv := trello.NewServer(trelloConfig)
		err := tSrv.Connect()
		if err != nil {
			log.Fatalf("Can't connect to trello: %s", err)
		}

		boards, err := tSrv.GetBoards()
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

		trelloConfig.Board = board[:trello.IDLength]
		err = tSrv.SetBoard(board[:trello.IDLength])
		if err != nil {
			log.Fatalf("Can't set trello board: %s", err)
		}

		viper.Set("trello.board", &trelloConfig.Board)

		lists, err := tSrv.GetLists()
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
		trelloConfig.Lists.Todo = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select doing list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Lists.Doing = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select done list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Lists.Done = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select review list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Lists.Review = choice

		_ = survey.AskOne(&survey.Select{
			Message: "Please select bucket list",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Lists.Bucket = choice

		viper.Set("trello.lists", &trelloConfig.Lists)

		labels, err := tSrv.GetLabels()
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
		trelloConfig.Labels.Jira = choice[:trello.IDLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Blocked label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Labels.Blocked = choice[:trello.IDLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Task label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Labels.Task = choice[:trello.IDLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Bug label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Labels.Bug = choice[:trello.IDLength]

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Story label",
			Options: keys,
		}, &choice)

		removeKeyFromSlice(keys, choice)
		trelloConfig.Labels.Story = choice[:trello.IDLength]

		viper.Set("trello.labels", &trelloConfig.Labels)

		members, err := tSrv.GetMembers()
		if err != nil {
			log.Fatalf("can't get members for %s board", board)
		}

		fmt.Printf("\nList of users from %s board:\n", board[trello.IDLength+2:])
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
			_ = slice // Fix linter staticcheck SA4006: this value of `slice` is never used

			break
		}
	}
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
