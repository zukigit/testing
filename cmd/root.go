package cmd

import (
	"os"
	"time"

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
	for _, testcase := range tcs {
		testcase.InfoLog("running")
		if testcase.Is_function_nil() {
			testcase.ErrorLog("testcase function is nil, skipping execution")
		} else {
			// start time
			startTime := time.Now()

			testcase.Set_status(testcase.Run_function())

			// total elasped time or duration of testcase
			duration := time.Since(startTime)

			testcase.Set_duration(duration)
		}
		// lib.InfoLog.Println("done running!")
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
