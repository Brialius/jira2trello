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
	Config
	cli   *trello.Client
	board *trello.Board
}

func NewClient(cfg *Config) *Client {
	return &Client{
		Config: *cfg,
	}
}

func (t *Client) Connect() error {
	t.cli = trello.NewClient(t.APIKey, t.Token)
	if len(t.Board) > 0 {
		if err := t.SetBoard(t.Board); err != nil {
			return err
		}
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
		res[board.ID] = &Board{
			URL:  board.URL,
			Name: board.Name,
			ID:   board.ID,
		}
	}

	t.writeToJSONFile(res, "boards.json")

	return res, nil
}

func (t *Client) GetLists() ([]*List, error) {
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

	t.writeToJSONFile(res, "lists.json")

	return res, nil
}

func (t *Client) GetLabels() ([]*Label, error) {
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

	t.writeToJSONFile(res, "labels.json")

	return res, nil
}

func (t *Client) GetBoardByID(id string) (*Board, error) {
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

func (t *Client) GetCards() ([]*Card, error) {
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
		})
	}

	t.writeToJSONFile(res, "cards.json")

	return res, nil
}

func (t *Client) writeToJSONFile(value interface{}, fileName string) {
	if t.Debug {
		b, _ := json.MarshalIndent(value, "", "  ")
		err := ioutil.WriteFile(fileName, b, 0644)

		if err != nil {
			fmt.Printf("can't write debug file: %s", err)
		}
	}
}

func (t *Client) CreateCard(card *Card) error {
	return t.cli.CreateCard(&trello.Card{
		Name:      card.Name,
		IDLabels:  *card.IDLabels,
		IDList:    card.ListID[:IDLength],
		IDMembers: strings.Split(card.IDMembers, ","),
		Desc:      card.Desc,
	}, trello.Defaults())
}

func (t *Client) MoveCardToList(cardID, listID string) error {
	card, err := t.cli.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.MoveToList(listID[:IDLength], trello.Defaults())
}

func (t *Client) UpdateCardLabels(cardID, labels string) error {
	card, err := t.cli.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.Update(trello.Arguments{"idLabels": labels})
}

func (t *Client) SetBoard(id string) error {
	t.Board = id
	board, err := t.cli.GetBoard(id, trello.Defaults())

	if err != nil {
		return err
	}

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
	return &t.Config
}
