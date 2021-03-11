/*
Copyright © 2021 Denis Belyatsky <denis.bel@gmail.com>

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
	"github.com/mattn/go-colorable"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

type SyncService struct {
	jCli   JiraConnector
	tCli   TrelloConnector
	jTasks map[string]*jira.Task
	tCards map[string]*trello.Card
}

func NewSyncService(jCli JiraConnector, tCli TrelloConnector) *SyncService {
	return &SyncService{
		jCli: jCli,
		tCli: tCli,
	}
}

func (s *SyncService) Sync() {
	var err error

	if err := s.jCli.Connect(); err != nil {
		log.Fatalf("Can't connect to jira server: %s", err)
	}

	if err := s.tCli.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	fmt.Print("Getting Jira tasks... ")

	if s.jTasks, err = s.jCli.GetUserTasks(); err != nil {
		log.Fatalf("can't get jira tasks: %s", err)
	}

	fmt.Printf("found %d\n", len(s.jTasks))

	fmt.Println()
	s.printJiraTasks(colorable.NewColorableStdout())
	fmt.Println()

	if s.tCards, err = getTrelloCards(s.tCli); err != nil {
		log.Fatalf("can't get trello cards: %s", err)
	}

	if err := s.syncTasks(); err != nil {
		log.Fatalf("can't sync tasks: %s", err)
	}

	if err := s.syncCompletedTasks(); err != nil {
		log.Fatalf("can't sync completed tasks: %s", err)
	}
}

func (s *SyncService) syncCompletedTasks() error {
	fmt.Println("Searching completed tasks..")

	for key, tCard := range s.tCards {
		if _, ok := s.jTasks[key]; !ok {
			if tCard.ListID != s.tCli.GetConfig().Lists.Done {
				if err := s.tCli.MoveCardToList(tCard.ID, s.tCli.GetConfig().Lists.Done); err != nil {
					return fmt.Errorf("can't move card to `Done` list: %w", err)
				}

				fmt.Printf("%s is completed!\n", key)
			}
		}
	}

	return nil
}

func (s *SyncService) syncTasks() error {
	fmt.Println("Sync tasks...")

	for key, jTask := range s.jTasks {
		listID := s.tCli.GetConfig().Lists.Todo
		labels := make([]string, 0)
		labels = append(labels, s.tCli.GetConfig().Labels.Jira)

		switch jTask.Status {
		case "In Progress", "In Dev / In Progress":
			listID = s.tCli.GetConfig().Lists.Doing
		case "Dependency", "Blocked":
			listID = s.tCli.GetConfig().Lists.Doing
			labels = append(labels, s.tCli.GetConfig().Labels.Blocked)
		case "jTask.Status", "In QA Review":
			listID = s.tCli.GetConfig().Lists.Review
		}

		switch jTask.Type {
		case "Story":
			labels = append(labels, s.tCli.GetConfig().Labels.Story)
		case "User Story":
			labels = append(labels, s.tCli.GetConfig().Labels.Story)
		case "Bug":
			labels = append(labels, s.tCli.GetConfig().Labels.Bug)
		default:
			labels = append(labels, s.tCli.GetConfig().Labels.Task)
		}

		if tCard, ok := s.tCards[key]; !ok {
			if err := s.addCardToList(jTask, listID, key, labels); err != nil {
				return fmt.Errorf("can't add task to list: %w", err)
			}
		} else {
			if err := s.updateCardLabels(tCard, labels); err != nil {
				return err
			}
			if err := s.updateCardList(tCard, listID, jTask); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SyncService) updateCardList(tCard *trello.Card, listID string, task *jira.Task) error {
	if tCard.ListID != listID {
		if listID == s.tCli.GetConfig().Lists.Doing || listID == s.tCli.GetConfig().Lists.Todo {
			if tCard.IsInAnyOfLists([]string{
				s.tCli.GetConfig().Lists.Bucket,
				s.tCli.GetConfig().Lists.Review,
			}) {
				return nil
			}
		}

		fmt.Printf("Moving %s to %s list\n", task.Key, trello.GetListNameByID(listID, s.tCli.GetConfig().Lists))
		err := s.tCli.MoveCardToList(tCard.ID, listID)

		if err != nil {
			return fmt.Errorf("can't move card to list: %w", err)
		}
	}

	return nil
}

func (s *SyncService) updateCardLabels(tCard *trello.Card, labels []string) error {
	if !reflect.DeepEqual(*tCard.IDLabels, labels) {
		fmt.Printf("Updating labels for %s\n", tCard.Key)
		err := s.tCli.UpdateCardLabels(tCard.ID, strings.Join(labels, ","))

		if err != nil {
			return fmt.Errorf("can't update labels on card `%s`: %w", tCard.Key, err)
		}
	}

	return nil
}

func (s *SyncService) addCardToList(task *jira.Task, listID string, key string, labels []string) error {
	fmt.Printf("Adding %s to %s list..\n", task.Key, trello.GetListNameByID(listID, s.tCli.GetConfig().Lists))
	desc := task.Desc + "\nJira link: " + task.Link + "\nType: " + task.Type

	if task.ParentKey != "" {
		desc += "\nParent link: " + task.ParentLink
	}

	return s.tCli.CreateCard(&trello.Card{
		Name:      key + " | " + task.Summary,
		ListID:    listID,
		Desc:      desc,
		IDLabels:  &labels,
		IDMembers: s.tCli.GetConfig().UserID,
	})
}

func getTrelloCards(tCli TrelloConnector) (map[string]*trello.Card, error) {
	fmt.Print("Getting Trello cards... ")

	tCards := map[string]*trello.Card{}

	cards, err := tCli.GetUserJiraCards()
	if err != nil {
		return nil, err
	}

	fmt.Printf("found %d\n", len(cards))

	for _, card := range cards {
		tCards[card.Key] = card
	}

	return tCards, nil
}

func (s *SyncService) printJiraTasks(out io.Writer) {
	w := new(tabwriter.Writer)

	w.Init(out, 0, 0, 4, ' ', tabwriter.FilterHTML+tabwriter.StripEscape)

	list := make([]*jira.Task, 0, len(s.jTasks))

	for _, task := range s.jTasks {
		list = append(list, task)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Created.Before(list[j].Created)
	})

	for _, task := range list {
		_, _ = fmt.Fprintln(w, task.TabString())
	}

	_ = w.Flush()
}
