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
	"github.com/Brialius/jira2trello/internal"
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/mattn/go-colorable"
	"io"
	"log"
	"sort"
)

const (
	orderDone = iota
	orderDoing
	orderReview
)

func Report(tCli TrelloConnector, jiraURL string) {
	if err := tCli.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	tCards, err := tCli.GetUserJiraCards()
	if err != nil {
		log.Fatalf("can't get trello cards: %s", err)
	}

	printReport(colorable.NewColorableStdout(), tCli, tCards, jiraURL)
}

func printReport(out io.Writer, tCli TrelloConnector, tCards []*trello.Card, jiraURL string) {
	var (
		done       int
		inProgress int
		inReview   int
	)

	sort.Slice(tCards, func(i, j int) bool {
		l := map[string]int{
			tCli.GetConfig().Lists.Done:   orderDone,
			tCli.GetConfig().Lists.Doing:  orderDoing,
			tCli.GetConfig().Lists.Review: orderReview,
		}

		return l[tCards[i].ListID] < l[tCards[j].ListID]
	})

	_, _ = fmt.Fprintln(out, "\n----------------------------------")

	doneString := internal.Green + "Done" + internal.ColorOff
	doingString := internal.Yellow + "In progress" + internal.ColorOff
	reviewString := internal.Cyan + "In review" + internal.ColorOff

	for _, tCard := range tCards {
		switch {
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Done}):
			printReportCard(out, tCard, doneString, jiraURL)
			done++
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Doing}):
			printReportCard(out, tCard, doingString, jiraURL)
			inProgress++
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Review}):
			printReportCard(out, tCard, reviewString, jiraURL)
			inReview++
		}
	}

	_, _ = fmt.Fprintln(out, "\n----------------------------------")
	_, _ = fmt.Fprintln(out, doingString+":", inProgress)
	_, _ = fmt.Fprintln(out, reviewString+":", inReview)
	_, _ = fmt.Fprintln(out, doneString+":", done)
}

func printReportCard(out io.Writer, tCard *trello.Card, status string, jiraURL string) {
	httpPrefix := internal.Blue + jiraURL + "/browse/"

	_, _ = fmt.Fprintf(out, "\n%s - %s\n", tCard.Name, status)
	_, _ = fmt.Fprintln(out, httpPrefix+tCard.Key+internal.ColorOff)
}
