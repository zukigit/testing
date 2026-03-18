package cmd

import (
	"fmt"
	"os"
	"strconv"

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

		testcase.StatusLog(string(status))
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
	Use:   "testing [ticket number]",
	Short: "run all tickets or a specific ticket",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			ticketNumStr := args[0]
			ticketNum, err := strconv.Atoi(ticketNumStr)
			if err != nil || ticketNum <= 0 {
				fmt.Println("Invalid ticket number. Must be a positive integer.")
				os.Exit(1)
			}

			var found bool
			for _, ticket := range ts {
				ticket.Prepare()
				if ticket.GetTicketNo() == uint(ticketNum) {
					tcs = ticket.Get_testcases()
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("Ticket %d not found\n", ticketNum)
				os.Exit(1)
			}
		} else {
			collectTcs()
		}
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
