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
	"github.com/Brialius/jira2trello/internal/jira"
	"io"
	"sort"
	"text/tabwriter"
)

func printJiraTasks(out io.Writer, jTasks map[string]*jira.Task) {
	const padding = 4

	textWriter := new(tabwriter.Writer)

	textWriter.Init(out, 0, 0, padding, ' ', tabwriter.FilterHTML+tabwriter.StripEscape)

	list := make([]*jira.Task, 0, len(jTasks))

	for _, task := range jTasks {
		list = append(list, task)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Created.Before(list[j].Created)
	})

	for _, task := range list {
		_, _ = fmt.Fprintln(textWriter, task.TabString())
	}

	_ = textWriter.Flush()
}
