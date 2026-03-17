package cmd

import (
	"fmt"
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
		// lib.InfoLog.Println("running")
		if testcase.Is_function_nil() {
			testcase.Set_status(models.FAILED)
			// lib.InfoLog.Printf("testcase: %d has no function. SKIPPED", testcase.Get_testcase_no())
		} else {
			// start time
			startTime := time.Now()

			testcase.Set_status(testcase.Run_function())

			// total elasped time or duration of testcase
			duration := time.Since(startTime)
			durationStr := fmt.Sprintf("%02d:%02d:%02d", int(duration/time.Hour), int(duration/time.Minute)%60, int(duration/time.Second)%60)

			testcase.Set_duration(durationStr)
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
