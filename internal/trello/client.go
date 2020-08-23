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
	"encoding/json"
	"fmt"
	"github.com/adlio/trello"
	"io/ioutil"
	"strings"
)

type Client struct {
	*Config
	cli   *trello.Client
	board *trello.Board
}

func NewClient(cfg *Config) *Client {
	return &Client{
		Config: cfg,
	}
}

func (t *Client) Connect() error {
	t.cli = trello.NewClient(t.APIKey, t.Token)
	if len(t.Board) > 0 {
		return t.SetBoard()
	}

	return nil
}

func (t *Client) GetBoards() (map[string]*Board, error) {
	res := map[string]*Board{}
	boards, err := t.cli.GetMyBoards(trello.Defaults())

	if err != nil {
		return nil, err
	}

	for _, board := range boards {
		res[board.Name] = &Board{
			URL:  board.URL,
			Name: board.Name,
			ID:   board.ID,
		}
	}

	t.writeToJSONFile(boards, "boards.json")
	t.writeToJSONFile(res, "boards_result.json")

	return res, nil
}

func (t *Client) GetLists() (map[string]*List, error) {
	lists, err := t.board.GetLists(trello.Defaults())
	if err != nil {
		return nil, err
	}

	res := map[string]*List{}
	for _, list := range lists {
		res[list.Name] = &List{
			Name: list.Name,
			ID:   list.ID,
		}
	}

	t.writeToJSONFile(lists, "lists.json")
	t.writeToJSONFile(res, "lists_result.json")

	return res, nil
}

func (t *Client) GetLabels() (map[string]*Label, error) {
	res := map[string]*Label{}

	labels, err := t.board.GetLabels(trello.Defaults())
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		res[label.Name] = &Label{
			Name: label.Name,
			ID:   label.ID,
		}
	}

	t.writeToJSONFile(labels, "labels.json")
	t.writeToJSONFile(res, "labels_result.json")

	return res, nil
}

func (t *Client) GetUserJiraCards() ([]*Card, error) {
	cards, err := t.board.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	res := make([]*Card, 0, len(cards))

	for _, card := range cards {
		if strings.Contains(strings.Join(card.IDMembers, ","), t.UserID) &&
			strings.Contains(strings.Join(card.IDLabels, ","), t.Labels.Jira) {
			res = append(res, &Card{
				ID:        card.ID,
				Name:      card.Name,
				ListID:    card.IDList,
				List:      GetListNameByID(card.IDList, t.Lists),
				Key:       strings.TrimSpace(strings.Split(card.Name, "|")[0]),
				Desc:      card.Desc,
				IDLabels:  &card.IDLabels,
				IDMembers: strings.Join(card.IDMembers, ","),
			})
		}
	}

	t.writeToJSONFile(cards, "cards.json")
	t.writeToJSONFile(res, "cards_result.json")

	return res, nil
}

func (t *Client) writeToJSONFile(value interface{}, fileName string) {
	if t.Debug {
		b, _ := json.MarshalIndent(value, "", "  ")
		err := ioutil.WriteFile(fileName, b, 0600)

		if err != nil {
			fmt.Printf("can't write debug file: %s", err)
		}
	}
}

func (t *Client) CreateCard(card *Card) error {
	return t.cli.CreateCard(&trello.Card{
		Name:      card.Name,
		IDLabels:  *card.IDLabels,
		IDList:    card.ListID,
		IDMembers: strings.Split(card.IDMembers, ","),
		Desc:      card.Desc,
	}, trello.Defaults())
}

func (t *Client) MoveCardToList(cardID, listID string) error {
	card, err := t.cli.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.MoveToList(listID, trello.Defaults())
}

func (t *Client) UpdateCardLabels(cardID, labels string) error {
	card, err := t.cli.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.Update(trello.Arguments{"idLabels": labels})
}

func (t *Client) SetBoard() error {
	board, err := t.cli.GetBoard(t.Board, trello.Defaults())
	if err != nil {
		return err
	}

	fmt.Printf("Using board: %s\n", board.Name)
	t.board = board

	return nil
}

func (t *Client) GetSelfMemberID() (string, error) {
	member, err := t.cli.GetMember("me", trello.Defaults())
	if err != nil {
		return "", err
	}

	t.writeToJSONFile(member, "self_id.json")

	return member.ID, nil
}

func (t *Client) GetConfig() *Config {
	return t.Config
}
