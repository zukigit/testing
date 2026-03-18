package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zukigit/testing/models"
	tickets "github.com/zukigit/testing/tickets"
)

var ts []models.Ticket
var tcs []models.TestCase
var failedTcs map[uint]uint

func collectTcs() {
	for _, ticket := range ts {
		ticket.Prepare()
		tcs = append(tcs, ticket.Get_testcases()...)
	}
}

func runTcs() {
	failedCount := 0
	failedTcs = make(map[uint]uint)

	for _, testcase := range tcs {
		testcase.InfoLog("running")
		status := testcase.Failed()
		if testcase.IsFunctionNil() {
			testcase.ErrorLog("testcase function is nil, skipping execution")
			status = testcase.Failed()
		} else {
			status = testcase.RunFunction()
		}

		if status != testcase.Passed() {
			failedCount++
			failedTcs[testcase.GetTicketNo()] = testcase.GetTestcaseNo()
		}

		testcase.InfoLog(string(status))
	}

	fmt.Println("@@@ FINISHED @@@")

	if failedCount > 0 {
		fmt.Println("Not Passed testcases:")
		for ticketNo, testcaseNo := range failedTcs {
			fmt.Printf("TicketNo: %d, TestcaseNo: %d\n", ticketNo, testcaseNo)
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
