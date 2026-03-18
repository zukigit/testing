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

func checkDuplicatesAndPrepare() {
	seen := make(map[uint]bool)
	var hasDuplicates bool
	for _, ticket := range ts {
		ticket.Prepare()
		no := ticket.GetTicketNo()
		if seen[no] {
			if !hasDuplicates {
				fmt.Println("Duplicated ticket numbers found:")
				hasDuplicates = true
			}
			fmt.Printf("- TicketNo: %d\n", no)
		}
		seen[no] = true
	}
	if hasDuplicates {
		os.Exit(1)
	}
}

func collectTcs() {
	for _, ticket := range ts {
		tcs = append(tcs, ticket.GetTestcases()...)
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

		testcase.StatusLog(status)
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
	Use:   "testing [ticket number] [testcase number]",
	Short: "run all tickets, a specific ticket, or a specific testcase",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		checkDuplicatesAndPrepare()

		if len(args) >= 1 {
			ticketNumStr := args[0]
			ticketNum, err := strconv.Atoi(ticketNumStr)
			if err != nil || ticketNum <= 0 {
				fmt.Println("Invalid ticket number. Must be a positive integer.")
				os.Exit(1)
			}

			testcaseNum := 0
			if len(args) == 2 {
				testcaseNumStr := args[1]
				testcaseNumParsed, err := strconv.Atoi(testcaseNumStr)
				if err != nil || testcaseNumParsed <= 0 {
					fmt.Println("Invalid testcase number. Must be a positive integer.")
					os.Exit(1)
				}
				testcaseNum = testcaseNumParsed
			}

			var foundTicket bool
			for _, ticket := range ts {
				if ticket.GetTicketNo() == uint(ticketNum) {
					foundTicket = true

					if testcaseNum > 0 {
						var foundTestcase bool
						for _, tc := range ticket.GetTestcases() {
							if tc.GetTestcaseNo() == uint(testcaseNum) {
								tcs = append(tcs, tc)
								foundTestcase = true
								break
							}
						}
						if !foundTestcase {
							fmt.Printf("Testcase %d not found in Ticket %d\n", testcaseNum, ticketNum)
							os.Exit(1)
						}
					} else {
						tcs = ticket.GetTestcases()
					}
					break
				}
			}

			if !foundTicket {
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
	ts = append(ts, new(tickets.Ticket2))
}
