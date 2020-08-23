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

		tCfg := trello.Config{}
		tCfg.Debug = Debug

		trelloQs := []*survey.Question{
			{
				Name: "apiKey",
				Prompt: &survey.Input{
					Help:    "API key can be generated here: https://trello.com/app-key",
					Message: "What is trello API key?",
					Default: viper.GetString("trello.apiKey"),
				},
				Validate: survey.Required,
			},
			{
				Name: "token",
				Prompt: &survey.Password{
					Help:    "Token can be generated here: https://trello.com/app-key",
					Message: "What is trello token?",
				},
			},
		}

		_ = survey.Ask(trelloQs, &tCfg)

		if tCfg.Token == "" {
			tCfg.Token = viper.GetString("trello.token")
		}

		viper.Set("trello.apiKey", tCfg.APIKey)
		viper.Set("trello.token", tCfg.Token)

		tCli := trello.NewClient(&tCfg)

		if err := tCli.Connect(); err != nil {
			log.Fatalf("Can't connect to trello: %s", err)
		}

		userID, err := tCli.GetSelfMemberID()
		if err != nil {
			log.Fatalf("can't get self id: %s", err)
		}

		tCfg.UserID = userID

		viper.Set("trello.userid", tCfg.UserID)

		boards, err := tCli.GetBoards()
		if err != nil {
			log.Fatalf("Can't get trello boards: %s", err)
		}

		boardNames := make([]string, 0, len(boards))

		for name := range boards {
			boardNames = append(boardNames, name)
		}

		sort.Strings(boardNames)

		var board string
		_ = survey.AskOne(&survey.Select{
			Message: "Please select trello board",
			Options: boardNames,
		}, &board)

		tCfg.Board = boards[board].ID

		if err := tCli.SetBoard(); err != nil {
			log.Fatalf("Can't set trello board: %s", err)
		}

		viper.Set("trello.board", &tCfg.Board)

		lists, err := tCli.GetLists()
		if err != nil {
			log.Fatalf("Can't get trello lists: %s", err)
		}

		listNames := make([]string, 0, len(lists))

		for name := range lists {
			listNames = append(listNames, name)
		}

		sort.Strings(listNames)

		var list string

		_ = survey.AskOne(&survey.Select{
			Message: "Please select todo list",
			Options: listNames,
		}, &list)

		removeKeyFromSlice(listNames, list)
		tCfg.Lists.Todo = lists[list].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select doing list",
			Options: listNames,
		}, &list)

		removeKeyFromSlice(listNames, list)
		tCfg.Lists.Doing = lists[list].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select done list",
			Options: listNames,
		}, &list)

		removeKeyFromSlice(boardNames, list)
		tCfg.Lists.Done = lists[list].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select review list",
			Options: listNames,
		}, &list)

		removeKeyFromSlice(listNames, list)
		tCfg.Lists.Review = lists[list].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select bucket list",
			Options: listNames,
		}, &list)

		tCfg.Lists.Bucket = lists[list].ID

		viper.Set("trello.lists", &tCfg.Lists)

		labels, err := tCli.GetLabels()
		if err != nil {
			log.Fatalf("Can't get trello labels: %s", err)
		}

		labelNames := make([]string, 0, len(labels))

		for name := range labels {
			labelNames = append(labelNames, name)
		}

		sort.Strings(labelNames)

		var label string

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Jira label",
			Options: labelNames,
		}, &label)

		removeKeyFromSlice(labelNames, label)
		tCfg.Labels.Jira = labels[label].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Blocked label",
			Options: labelNames,
		}, &label)

		removeKeyFromSlice(labelNames, label)
		tCfg.Labels.Blocked = labels[label].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Task label",
			Options: labelNames,
		}, &label)

		removeKeyFromSlice(labelNames, label)
		tCfg.Labels.Task = labels[label].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Bug label",
			Options: labelNames,
		}, &label)

		removeKeyFromSlice(labelNames, label)
		tCfg.Labels.Bug = labels[label].ID

		_ = survey.AskOne(&survey.Select{
			Message: "Please select Story label",
			Options: labelNames,
		}, &label)

		tCfg.Labels.Story = labels[label].ID

		viper.Set("trello.labels", &tCfg.Labels)

		tCfg.Debug = false

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
