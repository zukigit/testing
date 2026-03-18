package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zukigit/testing/models"
	tickets "github.com/zukigit/testing/tickets"
)

var ts []models.Ticket
var tcs []models.TestCase

func collectTcs() {
	for _, ticket := range ts {
		ticket.Prepare()
		tcs = append(tcs, ticket.Get_testcases()...)
	}
}

func runTcs() {
	failedCount := 0
	for _, testcase := range tcs {
		testcase.InfoLog("running")
		if testcase.IsFunctionNil() {
			testcase.ErrorLog("testcase function is nil, skipping execution")
			testcase.SetStatus(models.FAILED)
			failedCount++
		} else {
			status := testcase.Run_function()
			if status != testcase.Passed() {
				failedCount++
			}
			testcase.SetStatus(status)
		}

		testcase.InfoLog(string(testcase.GetStatus()))
	}

	if failedCount > 0 {
		for _, testcase := range tcs {
			if testcase.GetStatus() != testcase.Passed() {
				testcase.InfoLog(string(testcase.GetStatus()))
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "testing",
	Short: "run all tickets",
	Run: func(cmd *cobra.Command, args []string) {
		collectTcs()
		runTcs()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add your tickets here
	ts = append(ts, new(tickets.Ticket1))
}
