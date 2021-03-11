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
package cmd

import (
	"github.com/Brialius/jira2trello/internal/app"
	"github.com/Brialius/jira2trello/internal/jira"
	"github.com/spf13/viper"
	"log"

	"github.com/spf13/cobra"
)

// weeklyReportCmd represents the weekly-report command.
var weeklyReportCmd = &cobra.Command{
	Use:   "weekly-report",
	Short: "Weekly report based on jira query",
	Long:  "Get `Closed` or `In progress` stories and jira tasks with non zero `timespent` for last 7 days",
	Aliases: []string{
		"weekly",
	},
	Run: func(cmd *cobra.Command, args []string) {
		var jCfg jira.Config
		if err := viper.UnmarshalKey("jira", &jCfg); err != nil {
			log.Fatalf("Can't parse Jira config: %s", err)
		}

		jCfg.Debug = Debug

		app.WeeklyReport(jira.NewClient(&jCfg))
	},
}

func init() {
	rootCmd.AddCommand(weeklyReportCmd)
}
