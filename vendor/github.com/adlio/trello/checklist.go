// Copyright © 2016 Aaron Longwell
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package trello

import "fmt"

// Checklist represents Trello card's checklists.
// A card can have one zero or more checklists.
// https://developers.trello.com/reference/#checklist-object
type Checklist struct {
	client     *Client
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDBoard    string      `json:"idBoard,omitempty"`
	IDCard     string      `json:"idCard,omitempty"`
	Card       *Card       `json:"-"`
	Pos        float64     `json:"pos,omitempty"`
	CheckItems []CheckItem `json:"checkItems,omitempty"`
}

// CheckItem is a nested resource representing an item in Checklist.
type CheckItem struct {
	client      *Client
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	State       string     `json:"state"`
	IDChecklist string     `json:"idChecklist,omitempty"`
	Checklist   *Checklist `json:"-"`
	Pos         float64    `json:"pos,omitempty"`
}

// CheckItemState represents a CheckItem when it appears in CheckItemStates on a Card.
type CheckItemState struct {
	IDCheckItem string `json:"idCheckItem"`
	State       string `json:"state"`
}

// CreateChecklist creates a checklist.
// Attribute currently supported as extra argument: pos.
// Attributes currently known to be unsupported: idChecklistSource.
//
// API Docs: https://developers.trello.com/reference#cardsidchecklists-1
func (c *Client) CreateChecklist(card *Card, name string, extraArgs ...Arguments) (checklist *Checklist, err error) {
	path := "cards/" + card.ID + "/checklists"
	args := Arguments{
		"name": name,
		"pos":  "bottom",
	}

	args.flatten(extraArgs)

	checklist = &Checklist{}
	err = c.Post(path, args, &checklist)
	if err == nil {
		checklist.SetClient(c)
		checklist.IDCard = card.ID
		checklist.Card = card
		card.Checklists = append(card.Checklists, checklist)
	}
	return
}

// CreateCheckItem creates a checkitem inside the checklist.
// Attribute currently supported as extra argument: pos.
// Attributes currently known to be unsupported: checked.
//
// API Docs: https://developers.trello.com/reference#checklistsidcheckitems
func (cl *Checklist) CreateCheckItem(name string, extraArgs ...Arguments) (item *CheckItem, err error) {
	args := flattenArguments(extraArgs)
	return cl.client.CreateCheckItem(cl, name, args)
}

// CreateCheckItem creates a checkitem inside the given checklist.
// Attribute currently supported as extra argument: pos.
// Attributes currently known to be unsupported: checked.
//
// API Docs: https://developers.trello.com/reference#checklistsidcheckitems
func (c *Client) CreateCheckItem(checklist *Checklist, name string, extraArgs ...Arguments) (item *CheckItem, err error) {
	path := "checklists/" + checklist.ID + "/checkItems"
	args := Arguments{
		"name":    name,
		"pos":     "bottom",
		"checked": "false",
	}

	args.flatten(extraArgs)

	item = &CheckItem{}
	err = c.Post(path, args, item)
	if err == nil {
		checklist.CheckItems = append(checklist.CheckItems, *item)
	}
	return
}

// GetChecklist receives a checklist id and Arguments and returns the checklist if found
// with the credentials given for the receiver Client. Returns an error
// otherwise.
func (c *Client) GetChecklist(checklistID string, args Arguments) (checklist *Checklist, err error) {
	path := fmt.Sprintf("checklists/%s", checklistID)
	err = c.Get(path, args, &checklist)
	if checklist != nil {
		checklist.SetClient(c)
	}
	return checklist, err
}

// SetClient can be used to override this Checklist's internal connection to the
// Trello API. Normally, this is set automatically after API calls.
func (cl *Checklist) SetClient(newClient *Client) {
	cl.client = newClient
	for _, checkitem := range cl.CheckItems {
		checkitem.SetClient(newClient)
		checkitem.Checklist = cl
	}
}

// SetClient can be used to override this CheckItem's internal connection to the
// Trello API. Normally, this is set automatically after API calls.
func (ci *CheckItem) SetClient(newClient *Client) {
	ci.client = newClient
}
