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
	"github.com/Brialius/jira2trello/internal/trello"
	"github.com/spf13/viper"
	"log"

	"github.com/spf13/cobra"
)

var reportHTML bool
var reportWeekly bool

// reportCmd represents the report command.
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report based on trello cards",
	Long:  "Report based on trello cards",
	Run: func(cmd *cobra.Command, args []string) {
		var tCfg trello.Config
		var jCfg jira.Config

		if err := viper.UnmarshalKey("trello", &tCfg); err != nil {
			log.Fatalf("Can't parse Trello config: %s", err)
		}
		if err := viper.UnmarshalKey("jira", &jCfg); err != nil {
			log.Fatalf("Can't parse Jira config: %s", err)
		}

		tCfg.Debug = Debug

		app.Report(trello.NewClient(&tCfg), jira.NewClient(&jCfg), viper.GetString("jira.url"), reportHTML, reportWeekly)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().BoolVar(&reportHTML, "html", false,
		"generate html report and archive done cards")
	reportCmd.Flags().BoolVar(&reportWeekly, "weekly", false,
		"generate weekly report based on jira tasks")
}
