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
package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func Report() {
	tSrv := NewTrelloServer()
	err := tSrv.Connect()
	if err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	if err := viper.UnmarshalKey("users", &users); err != nil {
		log.Fatalf("Can't get user config: %s", err)
	}

	for _, user := range users {
		fmt.Printf("---------------------------------\n"+
			"User: %s\n"+
			"---------------------------------\n", user.Name)

		fmt.Println("Getting Trello cards...")
		trelloTasks := map[string]*TrelloCard{}
		cards, _ := tSrv.GetCards()
		for _, card := range cards {
			for _, labelId := range *card.IDLabels {
				if labelId == tSrv.Labels.Jira && strings.Contains(card.IDMembers, user.TrelloId) {
					trelloTasks[card.Key] = card
				}
			}
		}

		fmt.Println("Searching current tasks..")
		for _, tTask := range trelloTasks {
			if tTask.ListID == tSrv.Lists.Done[:TrelloIdLength] {
				fmt.Println(tTask.Name + " - Done")
				fmt.Println("https://jira.inbcu.com/browse/" + tTask.Key)
				fmt.Println("---------------------------------------------")
			} else if tTask.ListID == tSrv.Lists.Doing[:TrelloIdLength] {
				fmt.Println(tTask.Name + " - In progress")
				fmt.Println("https://jira.inbcu.com/browse/" + tTask.Key)
				fmt.Println("---------------------------------------------")
			}
		}
	}
}
