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
	"github.com/mattn/go-colorable"
	"html/template"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	doneString   = "Done"
	doingString  = "In progress"
	reviewString = "In review"
)

type Task struct {
	Name   string
	Status string
	Link   string
	Key    string
}

type report struct {
	HTMLReport   bool
	WeeklyReport bool
	Tasks        []*Task
	WeekNumber   int
	Year         int
}

func newReport(htmlReport bool, weeklyReport bool, tasks []*Task) *report {
	year, week := time.Now().ISOWeek()

	return &report{
		HTMLReport:   htmlReport,
		WeeklyReport: weeklyReport,
		Tasks:        tasks,
		WeekNumber:   week,
		Year:         year,
	}
}

// Generate report.
func (r *report) generate(out io.Writer) {
	if r.HTMLReport {
		t := template.Must(template.New("report").Parse(htmlTemplate))
		err := t.Execute(out, r)

		if err != nil {
			log.Fatalf("can't generate html report: %s", err)
		}

		return
	}

	_, _ = fmt.Fprintln(out, "\n----------------------------------")

	for _, t := range r.Tasks {
		_, _ = fmt.Fprintf(out, "\n%s\n", t)
	}

	_, _ = fmt.Fprintln(out, "\n----------------------------------")
}

func Report(tCli TrelloConnector, jCli JiraConnector, jiraURL string, reportHTML bool, reportWeekly bool) {
	var tasks []*Task

	if reportWeekly {
		tasks = WeeklyReportTasks(jCli)
	} else {
		tasks = trelloTasks(tCli, jiraURL)
	}

	r := newReport(reportHTML, reportWeekly, tasks)

	r.generate(r.getOutputWriter())

	// Archive done tasks if HTML report is generated.
	if reportHTML && !reportWeekly {
		if err := tCli.ArchiveAllCardsInList(tCli.GetConfig().Lists.Done); err != nil {
			log.Fatalf("can't archive done cards: %s", err)
		}

		fmt.Println("Done cards archived")
	}
}

// Determine destination writer
// depends on html report flag.
func (r *report) getOutputWriter() io.Writer {
	year := strconv.Itoa(r.Year)
	week := strconv.Itoa(r.WeekNumber)

	if r.HTMLReport {
		//nolint:gomnd,nosnakecase
		reportFile, err := os.OpenFile("jira2trello-report-"+year+"-"+week+".html",
			os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatalf("can't create report file: %s", err)
		}

		fmt.Printf("Report saved to %s\n", reportFile.Name())

		return reportFile
	}

	return colorable.NewColorableStdout()
}

func trelloTasks(tCli TrelloConnector, jiraURL string) []*Task {
	if err := tCli.Connect(); err != nil {
		log.Fatalf("Can't connect to trello: %s", err)
	}

	tCards, err := tCli.GetUserJiraCards()
	if err != nil {
		log.Fatalf("can't get trello cards: %s", err)
	}

	done := make([]*Task, 0)
	inProgress := make([]*Task, 0)
	inReview := make([]*Task, 0)

	sort.Slice(tCards, func(i, j int) bool {
		return tCards[i].Key < tCards[j].Key
	})

	tasks := make([]*Task, 0)

	for _, tCard := range tCards {
		switch {
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Done}):
			done = append(done, &Task{
				Name:   strings.TrimPrefix(tCard.Name, tCard.Key+" | "),
				Status: doneString,
				Link:   jiraURL + "/browse/" + tCard.Key,
				Key:    tCard.Key,
			})
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Doing}):
			inProgress = append(inProgress, &Task{
				Name:   strings.TrimPrefix(tCard.Name, tCard.Key+" | "),
				Status: doingString,
				Link:   jiraURL + "/browse/" + tCard.Key,
				Key:    tCard.Key,
			})
		case tCard.IsInAnyOfLists([]string{tCli.GetConfig().Lists.Review}):
			inReview = append(inReview, &Task{
				Name:   strings.TrimPrefix(tCard.Name, tCard.Key+" | "),
				Status: reviewString,
				Link:   jiraURL + "/browse/" + tCard.Key,
				Key:    tCard.Key,
			})
		}
	}

	tasks = append(tasks, done...)
	tasks = append(tasks, inProgress...)
	tasks = append(tasks, inReview...)

	return tasks
}

func (t *Task) String() string {
	return fmt.Sprintf("%s | %s - %s\n%s", t.Key, t.Name, t.Status, t.Link)
}
