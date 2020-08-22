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
package app

import (
	"fmt"
	"github.com/Brialius/jira2trello/internal/jira"
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

var users []*UserConfig

func Sync(jSrv *jira.Client, tSrv *trello.Client) {
	if err := jSrv.Connect(); err != nil {
		log.Fatalf("Can't connect to jira server: %s", err)
	}

	err := tSrv.Connect()
	if err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	if err := viper.UnmarshalKey("users", &users); err != nil {
		log.Fatalf("Can't get user config: %s", err)
	}

	for _, user := range users {
		fmt.Printf("---------------------------------\n"+
			"User: %s\n"+
			"---------------------------------\n", user.Name)

		fmt.Println("Getting Jira tasks...")
		jTasks, err := jSrv.GetUserTasks(user.Email)
		if err != nil {
			log.Fatalf("can't get tasks for %s from jira server: %s", user.Email, err)
		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)

		for _, task := range jTasks {
			_, _ = fmt.Fprintln(w, task.TabString())
		}
		_ = w.Flush()

		fmt.Println("Getting Trello cards...")
		tCards := map[string]*trello.Card{}
		cards, _ := tSrv.GetCards()
		for _, card := range cards {
			for _, labelId := range *card.IDLabels {
				if labelId == tSrv.Labels.Jira && strings.Contains(card.IDMembers, user.TrelloId) {
					tCards[card.Key] = card
				}
			}
		}

		fmt.Println("Sync tasks...")
		for key, value := range jTasks {
			list := tSrv.Lists.Todo
			labels := make([]string, 0)
			labels = append(labels, tSrv.Labels.Jira)
			switch {
			case strings.Contains(value.Status, "In Progress"):
				list = tSrv.Lists.Doing
			case strings.Contains(value.Status, "Dependency") || strings.Contains(value.Status, "Blocked"):
				list = tSrv.Lists.Doing
				labels = append(labels, tSrv.Labels.Blocked)
			}
			switch value.Type {
			case "Story":
				labels = append(labels, tSrv.Labels.Story)
			case "User Story":
				labels = append(labels, tSrv.Labels.Story)
			case "Bug":
				labels = append(labels, tSrv.Labels.Bug)
			default:
				labels = append(labels, tSrv.Labels.Task)
			}
			if tCard, ok := tCards[key]; !ok {
				fmt.Printf("Adding %s to %s list..\n", value.Key, list[trello.IdLength+3:])
				desc := value.Desc + "\nJira link: " + value.Link + "\nType: " + value.Type
				if value.ParentKey != "" {
					desc += "\nParent link: " + value.ParentLink
				}
				err = tSrv.CreateCard(&trello.Card{
					Name:      key + " | " + value.Summary,
					ListID:    list[:trello.IdLength],
					Desc:      desc,
					IDLabels:  &labels,
					IDMembers: user.TrelloId,
				})
				if err != nil {
					log.Fatalf("can't add task to list: %s", err)
				}
			} else {
				// Update labels
				if !reflect.DeepEqual(*tCard.IDLabels, labels) {
					fmt.Printf("Updating labels for %s\n", tCard.Key)
					err := tSrv.UpdateCardLabels(tCard.ID, strings.Join(labels, ","))
					if err != nil {
						log.Fatalf("can't update labels oncard `%s`: %s", tCard.Key, err)
					}
				}
				// Update trello card list
				if tCard.ListID != list[:trello.IdLength] {
					if list == tSrv.Lists.Doing || list == tSrv.Lists.Todo {
						if tCard.ListID == tSrv.Lists.Review[:trello.IdLength] || tCard.ListID == tSrv.Lists.Bucket[:trello.IdLength] {
							continue
						}
					}
					fmt.Printf("Moving %s to %s list\n", value.Key, list[trello.IdLength+3:])
					err = tSrv.MoveCardToList(tCard.ID, list[:trello.IdLength])
					if err != nil {
						log.Fatalf("can't move card to list: %s", err)
					}
				}
			}
		}

		fmt.Println("Searching completed tasks..")
		for key, tCard := range tCards {
			if _, ok := jTasks[key]; !ok {
				if tCard.ListID != tSrv.Lists.Done[:trello.IdLength] {
					err = tSrv.MoveCardToList(tCard.ID, tSrv.Lists.Done[:trello.IdLength])
					if err != nil {
						log.Fatalf("can't move card to `Done` list: %s", err)
					}
					fmt.Printf("%s is completed!\n", key)
				}
			}
		}
	}
}
