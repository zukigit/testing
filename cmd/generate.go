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

func (t *Ticket{{.TicketNum}}) New_testcase(testcase_id uint, testcase_description string) *models.TestCase {
	return models.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket{{.TicketNum}}) Get_ticket_no() uint {
	return t.Ticket_no
}

func (t *Ticket{{.TicketNum}}) Set_ticket_no(ticket_no uint) {
	t.Ticket_no = ticket_no
}

func (t *Ticket{{.TicketNum}}) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket{{.TicketNum}}) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket{{.TicketNum}}) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket{{.TicketNum}}) Get_ticket_description() string {
	return t.Ticket_description
}

func (t *Ticket{{.TicketNum}}) Set_ticket_description(testcase_description string) {
	t.Ticket_description = testcase_description
}

func (t *Ticket{{.TicketNum}}) Add_testcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket{{.TicketNum}}) Get_testcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket{{.TicketNum}}) Prepare() {
	t.Set_ticket_no({{.TicketNum}})
	t.Set_ticket_description("Enter your ticket description here.")

	// TESTCASE 1
	tc := t.New_testcase(1, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {
		//Enter your testcase function here
		return tc.Failed() // or tc.Passed() or tc.MustCheck()
	}
	tc.Set_function(tc_func)
	t.Add_testcase(tc)
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
