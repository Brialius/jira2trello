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
package trello

import (
	"fmt"
	"github.com/adlio/trello"
	"strings"
)

const IdLength = 24

type Server struct {
	Config
	cli   *trello.Client
	board *trello.Board
}

type Lists struct {
	Todo   string
	Doing  string
	Done   string
	Review string
	Bucket string
}

type Labels struct {
	Jira    string
	Blocked string
	Task    string
	Bug     string
	Story   string
}

type Card struct {
	ID        string
	Name      string
	ListID    string
	List      string
	Labels    string
	Key       string
	Desc      string
	IDLabels  *[]string
	IDMembers string
	cli       *trello.Client
}

type Board struct {
	URL  string
	Name string
	ID   string
}

type Label struct {
	Name string
	ID   string
}

type List struct {
	Name string
	ID   string
}

type Member struct {
	Name     string
	FullName string
	ID       string
}

func (c *Card) String() string {
	return fmt.Sprintf("%s | %s(%s): %s - %s", c.Key, c.List, c.ListID, c.Name, *c.IDLabels)
}

func NewServer(cfg Config) *Server {
	return &Server{
		Config: Config{
			ApiKey: cfg.ApiKey,
			Token:  cfg.Token,
			Board:  cfg.Board,
			Lists: Lists{
				Todo:   cfg.Lists.Todo,
				Doing:  cfg.Lists.Doing,
				Done:   cfg.Lists.Done,
				Review: cfg.Lists.Review,
				Bucket: cfg.Lists.Bucket,
			},
			Labels: Labels{
				Jira:    cfg.Labels.Jira,
				Blocked: cfg.Labels.Blocked,
				Bug:     cfg.Labels.Bug,
				Task:    cfg.Labels.Task,
				Story:   cfg.Labels.Story,
			},
		},
	}
}

func (t *Server) Connect() error {
	t.cli = trello.NewClient(t.ApiKey, t.Token)
	if len(t.Board) > 0 {
		if err := t.SetBoard(t.Board); err != nil {
			return err
		}
	}
	return nil
}

func (t *Server) GetBoards() (map[string]*Board, error) {
	res := map[string]*Board{}
	boards, err := t.cli.GetMyBoards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	for _, board := range boards {
		res[board.ID] = &Board{
			URL:  board.URL,
			Name: board.Name,
			ID:   board.ID,
		}
	}
	return res, nil
}

func (t *Server) GetLists() ([]*List, error) {
	lists, err := t.board.GetLists(trello.Defaults())
	if err != nil {
		return nil, err
	}

	res := make([]*List, 0, len(lists))
	for _, list := range lists {
		res = append(res, &List{
			Name: list.Name,
			ID:   list.ID,
		})
	}
	return res, nil
}

func (t *Server) GetLabels() ([]*Label, error) {
	res := make([]*Label, 0)
	board, err := t.cli.GetBoard(t.Board, trello.Defaults())
	if err != nil {
		return nil, err
	}
	labels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		res = append(res, &Label{
			Name: label.Name,
			ID:   label.ID,
		})
	}
	return res, nil
}

func (t *Server) GetBoardById(id string) (*Board, error) {
	board, err := t.cli.GetBoard(id, trello.Defaults())
	if err != nil {
		return nil, err
	}
	return &Board{
		URL:  board.URL,
		Name: board.Name,
		ID:   board.ID,
	}, err
}

func (t *Server) GetMembers() ([]*trello.Member, error) {
	members, err := t.board.GetMembers(trello.Defaults())
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (t *Server) GetCards() ([]*Card, error) {
	cards, err := t.board.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	res := make([]*Card, 0, len(cards))
	for _, card := range cards {
		res = append(res, &Card{
			ID:        card.ID,
			Name:      card.Name,
			ListID:    card.IDList,
			Key:       strings.TrimSpace(strings.Split(card.Name, "|")[0]),
			Desc:      card.Desc,
			IDLabels:  &card.IDLabels,
			IDMembers: strings.Join(card.IDMembers, ","),
			cli:       t.cli,
		})
	}
	return res, nil
}

func (t *Server) CreateCard(card *Card) error {
	return t.cli.CreateCard(&trello.Card{
		Name:      card.Name,
		IDLabels:  *card.IDLabels,
		IDList:    card.ListID,
		IDMembers: strings.Split(card.IDMembers, ","),
		Desc:      card.Desc,
	}, trello.Defaults())
}

func (c *Card) MoveToList(listId string) error {
	card, err := c.cli.GetCard(c.ID, trello.Defaults())
	if err != nil {
		return err
	}
	return card.MoveToList(listId, trello.Defaults())
}

func (c *Card) UpdateLabels(labels string) error {
	return c.updateArgs(trello.Arguments{"idLabels": labels})
}

func (c *Card) updateArgs(args trello.Arguments) error {
	card, err := c.cli.GetCard(c.ID, trello.Defaults())
	if err != nil {
		return err
	}
	return card.Update(args)
}

func (t *Server) SetBoard(id string) error {
	t.Board = id
	board, err := t.cli.GetBoard(id, trello.Defaults())
	if err != nil {
		return err
	}
	t.board = board
	return nil
}
