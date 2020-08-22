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
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

func Sync(jSrv *jira.Client, tSrv *trello.Client, users []*UserConfig) {
	if err := jSrv.Connect(); err != nil {
		log.Fatalf("Can't connect to jira server: %s", err)
	}

	if err := tSrv.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	for _, user := range users {
		fmt.Printf("---------------------------------\n"+
			"User: %s\n"+
			"---------------------------------\n", user.Name)

		jTasks, err := getJiraTasks(jSrv, user)
		if err != nil {
			log.Fatalf("can't get jira tasks: %s", err)
		}

		printJiraTasks(jTasks)

		tCards, err := getTrelloCards(tSrv, user)
		if err != nil {
			log.Fatalf("can't get trello cards: %s", err)
		}

		if err := syncTasks(jTasks, tSrv, tCards, user); err != nil {
			log.Fatalf("can't sync tasks: %s", err)
		}

		if err := syncCompletedTasks(tCards, jTasks, tSrv); err != nil {
			log.Fatalf("can't sync completed tasks: %s", err)
		}
	}
}

func syncCompletedTasks(tCards map[string]*trello.Card, jTasks map[string]*jira.Task, tSrv *trello.Client) error {
	fmt.Println("Searching completed tasks..")

	for key, tCard := range tCards {
		if _, ok := jTasks[key]; !ok {
			if tCard.ListID != tSrv.Lists.Done[:trello.IDLength] {
				if err := tSrv.MoveCardToList(tCard.ID, tSrv.Lists.Done[:trello.IDLength]); err != nil {
					return fmt.Errorf("can't move card to `Done` list: %w", err)
				}

				fmt.Printf("%s is completed!\n", key)
			}
		}
	}

	return nil
}

func syncTasks(jTasks map[string]*jira.Task, tSrv *trello.Client,
	tCards map[string]*trello.Card, user *UserConfig) error {
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
			if err := addCardToList(value, list, tSrv, key, labels, user); err != nil {
				return fmt.Errorf("can't add task to list: %w", err)
			}
		} else {
			if err := updateCardLabels(tCard, labels, tSrv); err != nil {
				return err
			}
			if err := updateCardList(tCard, list, tSrv, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func updateCardList(tCard *trello.Card, list string, tSrv *trello.Client, task *jira.Task) error {
	if tCard.ListID != list[:trello.IDLength] {
		if list == tSrv.Lists.Doing || list == tSrv.Lists.Todo {
			if tCard.IsInAnyOfLists([]string{
				tSrv.Lists.Bucket,
				tSrv.Lists.Review,
			}) {
				return nil
			}
		}

		fmt.Printf("Moving %s to %s list\n", task.Key, list[trello.IDLength+3:])
		err := tSrv.MoveCardToList(tCard.ID, list)

		if err != nil {
			return fmt.Errorf("can't move card to list: %w", err)
		}
	}

	return nil
}

func updateCardLabels(tCard *trello.Card, labels []string, tSrv *trello.Client) error {
	if !reflect.DeepEqual(*tCard.IDLabels, labels) {
		fmt.Printf("Updating labels for %s\n", tCard.Key)
		err := tSrv.UpdateCardLabels(tCard.ID, strings.Join(labels, ","))

		if err != nil {
			return fmt.Errorf("can't update labels on card `%s`: %w", tCard.Key, err)
		}
	}

	return nil
}

func addCardToList(task *jira.Task, list string, tSrv *trello.Client,
	key string, labels []string, user *UserConfig) error {
	fmt.Printf("Adding %s to %s list..\n", task.Key, list[trello.IDLength+3:])
	desc := task.Desc + "\nJira link: " + task.Link + "\nType: " + task.Type

	if task.ParentKey != "" {
		desc += "\nParent link: " + task.ParentLink
	}

	return tSrv.CreateCard(&trello.Card{
		Name:      key + " | " + task.Summary,
		ListID:    list,
		Desc:      desc,
		IDLabels:  &labels,
		IDMembers: user.TrelloID,
	})
}

func getTrelloCards(tSrv *trello.Client, user *UserConfig) (map[string]*trello.Card, error) {
	fmt.Println("Getting Trello cards...")

	tCards := map[string]*trello.Card{}

	cards, err := tSrv.GetCards()
	if err != nil {
		return nil, err
	}

	for _, card := range cards {
		for _, labelID := range *card.IDLabels {
			if labelID == tSrv.Labels.Jira && strings.Contains(card.IDMembers, user.TrelloID) {
				tCards[card.Key] = card
			}
		}
	}

	return tCards, nil
}

func printJiraTasks(jTasks map[string]*jira.Task) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)

	for _, task := range jTasks {
		_, _ = fmt.Fprintln(w, task.TabString())
	}

	_ = w.Flush()
}

func getJiraTasks(jSrv *jira.Client, user *UserConfig) (map[string]*jira.Task, error) {
	fmt.Println("Getting Jira tasks...")

	jTasks, err := jSrv.GetUserTasks(user.Email)

	if err != nil {
		log.Fatalf("can't get tasks for %s from jira server: %s", user.Email, err)
	}

	return jTasks, err
}
