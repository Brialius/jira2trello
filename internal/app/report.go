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
	"log"
	"sort"
)

func Report(tCli trello.Connector) {
	if err := tCli.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	tCards, err := tCli.GetUserJiraCards()
	if err != nil {
		log.Fatalf("can't get trello cards: %s", err)
	}

	sort.Slice(tCards, func(i, j int) bool {
		return tCards[i].List > tCards[j].List
	})

	var (
		done       int
		inProgress int
		inReview   int
	)

	fmt.Println("\n----------------------------------")

	for _, tTask := range tCards {
		switch {
		case tTask.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Done}):
			fmt.Printf("%s - Done\n", tTask.Name)
			fmt.Printf("https://jira.inbcu.com/browse/%s\n", tTask.Key)
			done++
		case tTask.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Doing}):
			fmt.Printf("%s - In progress\n", tTask.Name)
			fmt.Printf("https://jira.inbcu.com/browse/%s\n", tTask.Key)
			inProgress++
		case tTask.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Review}):
			fmt.Printf("%s - In review\n", tTask.Name)
			fmt.Printf("https://jira.inbcu.com/browse/%s\n", tTask.Key)
			inReview++
		}
	}

	fmt.Println("\n----------------------------------")
	fmt.Printf("In progress: %d\n", inProgress)
	fmt.Printf("In review: %d\n", inReview)
	fmt.Printf("Done: %d\n", done)
}
