package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zukigit/testing/models"
	tickets "github.com/zukigit/testing/tickets"
)

var ts []models.Ticket
var failedTcs map[uint][]uint
var skippedTickets []uint

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

// runTestcase executes a single testcase, records any failure, and returns the status.
func runTestcase(testcase models.TestCase) models.TestcaseStatus {
	testcase.InfoLog("running")

	var status models.TestcaseStatus
	if testcase.IsFunctionNil() {
		testcase.ErrorLog("testcase function is nil, skipping execution")
		status = testcase.Failed()
	} else {
		status = testcase.RunFunction()
	}

	if status != testcase.Passed() {
		ticketNo := testcase.GetTicketNo()
		failedTcs[ticketNo] = append(failedTcs[ticketNo], testcase.GetTestcaseNo())
	}

	testcase.StatusLog(status)
	return status
}

// runTicket runs all (or a specific) testcase(s) for a ticket.
// It creates a fresh context for the ticket, cancels it when done.
// Testcase 0 is treated as a preparation step and always runs first.
// If TC0 returns FAILED, the entire ticket is skipped and logged.
// filterTcNo == 0 means run all testcases (except testcase 0, which already ran).
func runTicket(ticket models.Ticket, filterTcNo uint) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticket.SetContext(ctx)

	// always run testcase 0 first as preparation, if it exists
	for _, testcase := range ticket.GetTestcases() {
		if testcase.GetTestcaseNo() == 0 {
			testcase.InfoLog("running preparation")
			if status := runTestcase(testcase); status == testcase.Failed() {
				testcase.ErrorLog("preparation failed — skipping entire ticket %d", ticket.GetTicketNo())
				skippedTickets = append(skippedTickets, ticket.GetTicketNo())
				return
			}
			break
		}
	}

	// run the remaining testcases (skip testcase 0, it already ran)
	for _, testcase := range ticket.GetTestcases() {
		if testcase.GetTestcaseNo() == 0 {
			continue
		}
		if filterTcNo > 0 && testcase.GetTestcaseNo() != filterTcNo {
			continue
		}
		runTestcase(testcase)
	}
}

func runAllTickets() {
	for _, ticket := range ts {
		runTicket(ticket, 0)
	}
}

func printResults() {
	fmt.Println("@@@ FINISHED @@@")
	if len(skippedTickets) > 0 {
		fmt.Println("Skipped tickets (preparation failed):")
		for _, ticketNo := range skippedTickets {
			fmt.Printf("- TicketNo: %d\n", ticketNo)
		}
	}
	if len(failedTcs) > 0 {
		fmt.Println("Not Passed testcases:")
		for ticketNo, testcaseNos := range failedTcs {
			for _, testcaseNo := range testcaseNos {
				fmt.Printf("TicketNo: %d, TestcaseNo: %d\n", ticketNo, testcaseNo)
			}
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "testing [ticket number] [testcase number]",
	Short: "run all tickets, a specific ticket, or a specific testcase",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		checkDuplicatesAndPrepare()
		failedTcs = make(map[uint][]uint)
		skippedTickets = nil

		if len(args) >= 1 {
			ticketNum, err := strconv.Atoi(args[0])
			if err != nil || ticketNum <= 0 {
				fmt.Println("Invalid ticket number. Must be a positive integer.")
				os.Exit(1)
			}

			testcaseNum := 0
			if len(args) == 2 {
				testcaseNumParsed, err := strconv.Atoi(args[1])
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
						// validate that the testcase exists before running
						var foundTestcase bool
						for _, tc := range ticket.GetTestcases() {
							if tc.GetTestcaseNo() == uint(testcaseNum) {
								foundTestcase = true
								break
							}
						}
						if !foundTestcase {
							fmt.Printf("Testcase %d not found in Ticket %d\n", testcaseNum, ticketNum)
							os.Exit(1)
						}
					}

					runTicket(ticket, uint(testcaseNum))
					break
				}
			}

			if !foundTicket {
				fmt.Printf("Ticket %d not found\n", ticketNum)
				os.Exit(1)
			}
		} else {
			runAllTickets()
		}

		printResults()
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
