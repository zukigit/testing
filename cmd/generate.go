package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [ticket number]",
	Short: "Generate a new ticket file",
	Long:  `Generate a new ticket file under tickets/ directory with the format of tickets/ticket_1.go.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketNumStr := args[0]
		ticketNum, err := strconv.Atoi(ticketNumStr)
		if err != nil || ticketNum <= 0 {
			fmt.Println("Invalid ticket number. Must be a positive integer.")
			os.Exit(1)
		}

		generateTicket(ticketNum)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

const ticketTemplate = `package tickets

import "github.com/zukigit/testing/models"

type Ticket{{.TicketNum}} struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []models.TestCase
}

func (t *Ticket{{.TicketNum}}) NewTestcase(testcaseNo uint, testcaseDescription string) *models.TestCase {
	return models.NewTestcase(t.Ticket_no, testcaseNo, testcaseDescription)
}

func (t *Ticket{{.TicketNum}}) GetTicketNo() uint {
	return t.Ticket_no
}

func (t *Ticket{{.TicketNum}}) SetTicketNo(ticketNo uint) {
	t.Ticket_no = ticketNo
}

func (t *Ticket{{.TicketNum}}) GetTicketDescription() string {
	return t.Ticket_description
}

func (t *Ticket{{.TicketNum}}) SetTicketDescription(testcaseDescription string) {
	t.Ticket_description = testcaseDescription
}

func (t *Ticket{{.TicketNum}}) AddTestcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket{{.TicketNum}}) GetTestcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket{{.TicketNum}}) Prepare() {
	t.SetTicketNo({{.TicketNum}})
	t.SetTicketDescription("Enter your ticket description here.")

	// TESTCASE 1
	tc := t.NewTestcase(1, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {
		//Enter your testcase function here

		//You can log as follow
		tc.ErrorLog("it is just example")
		return tc.Failed() // or tc.Passed() or tc.MustCheck()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)
}
`

func generateTicket(num int) {
	filename := filepath.Join("tickets", fmt.Sprintf("ticket_%d.go", num))
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File %s already exists\n", filename)
		os.Exit(1)
	}

	tmpl, err := template.New("ticket").Parse(ticketTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer f.Close()

	data := struct {
		TicketNum int
	}{
		TicketNum: num,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", filename)
}
