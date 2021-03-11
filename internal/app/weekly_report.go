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
	"github.com/mattn/go-colorable"
	"log"
)

func WeeklyReport(jCli JiraConnector) {
	if err := jCli.Connect(); err != nil {
		log.Fatalf("Can't connect to jira server: %s", err)
	}

	tasks, err := jCli.GetUserTasks("(status changed to closed after -7d  OR status = \"In Dev / In Progress\") AND (timespent != 0 OR issuetype = Story) ORDER BY priority DESC, updated DESC")
	if err != nil {
		log.Fatalf("Can't get jira tasks: %s", err)
	}

	printJiraTasks(colorable.NewColorableStdout(), tasks)
}