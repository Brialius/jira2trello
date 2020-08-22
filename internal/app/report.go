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
)

func Report(tSrv *trello.Client, users []*UserConfig) {
	err := tSrv.Connect()
	if err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	for _, user := range users {
		fmt.Printf("---------------------------------\n"+
			"User: %s\n"+
			"---------------------------------\n", user.Name)

		tCards := getTrelloTasks(tSrv, user)

		fmt.Println("Searching current trello tasks..")

		for _, tTask := range tCards {
			switch {
			case tTask.IsInAnyOfLists([]string{tSrv.Lists.Done}):
				fmt.Println(tTask.Name + " - Done")
				fmt.Println("https://jira.inbcu.com/browse/" + tTask.Key)
				fmt.Println("---------------------------------------------")
			case tTask.IsInAnyOfLists([]string{tSrv.Lists.Doing}):
				fmt.Println(tTask.Name + " - In progress")
				fmt.Println("https://jira.inbcu.com/browse/" + tTask.Key)
				fmt.Println("---------------------------------------------")
			case tTask.IsInAnyOfLists([]string{tSrv.Lists.Review}):
				fmt.Println(tTask.Name + " - In review")
				fmt.Println("https://jira.inbcu.com/browse/" + tTask.Key)
				fmt.Println("---------------------------------------------")
			}
		}
	}
}
