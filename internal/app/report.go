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
	"github.com/Brialius/jira2trello/internal/trello"
	"io"
	"log"
	"os"
	"sort"
)

func Report(tCli TrelloConnector) {
	if err := tCli.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	tCards, err := tCli.GetUserJiraCards()
	if err != nil {
		log.Fatalf("can't get trello cards: %s", err)
	}

	printReport(os.Stdout, tCli, tCards)
}

func printReport(out io.Writer, tCli TrelloConnector, tCards []*trello.Card) {
	var (
		done       int
		inProgress int
		inReview   int
	)

	sort.Slice(tCards, func(i, j int) bool {
		l := map[string]int{
			tCli.GetConfig().Lists.Done:   0,
			tCli.GetConfig().Lists.Doing:  1,
			tCli.GetConfig().Lists.Review: 2,
		}
		return l[tCards[i].ListID] < l[tCards[j].ListID]
	})

	_, _ = fmt.Fprintln(out, "\n----------------------------------")

	for _, tCard := range tCards {
		switch {
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Done}):
			_, _ = fmt.Fprintf(out, "\n%s - Done\n", tCard.Name)
			_, _ = fmt.Fprintf(out, "https://jira.inbcu.com/browse/%s\n", tCard.Key)
			done++
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Doing}):
			_, _ = fmt.Fprintf(out, "\n%s - In progress\n", tCard.Name)
			_, _ = fmt.Fprintf(out, "https://jira.inbcu.com/browse/%s\n", tCard.Key)
			inProgress++
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Review}):
			_, _ = fmt.Fprintf(out, "\n%s - In review\n", tCard.Name)
			_, _ = fmt.Fprintf(out, "https://jira.inbcu.com/browse/%s\n", tCard.Key)
			inReview++
		}
	}

	_, _ = fmt.Fprintln(out, "\n----------------------------------")
	_, _ = fmt.Fprintf(out, "In progress: %d\n", inProgress)
	_, _ = fmt.Fprintf(out, "In review: %d\n", inReview)
	_, _ = fmt.Fprintf(out, "Done: %d\n", done)
}
