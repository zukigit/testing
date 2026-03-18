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

import (
	"context"

	"github.com/zukigit/testing/models"
	"github.com/zukigit/testing/zabbix"
)

type Ticket{{.TicketNum}} struct {
	TicketNo          uint
	TicketDescription string
	Testcases         []models.TestCase
	context           context.Context

	// share objects
	zabbix zabbix.Zabbix
}

func (t *Ticket{{.TicketNum}}) SetZabbix(zabbix zabbix.Zabbix) {
	t.zabbix = zabbix
}

func (t *Ticket{{.TicketNum}}) GetZabbix() zabbix.Zabbix {
	return t.zabbix
}


func (t *Ticket{{.TicketNum}}) SetContext(ctx context.Context) {
	t.context = ctx
}

func (t *Ticket{{.TicketNum}}) GetContext() context.Context {
	return t.context
}

func (t *Ticket{{.TicketNum}}) NewTestcase(testcaseNo uint, testcaseDescription string) *models.TestCase {
	return models.NewTestcase(t.TicketNo, testcaseNo, testcaseDescription)
}

func (t *Ticket{{.TicketNum}}) GetTicketNo() uint {
	return t.TicketNo
}

func (t *Ticket{{.TicketNum}}) SetTicketNo(ticketNo uint) {
	t.TicketNo = ticketNo
}

func (t *Ticket{{.TicketNum}}) GetTicketDescription() string {
	return t.TicketDescription
}

func (t *Ticket{{.TicketNum}}) SetTicketDescription(testcaseDescription string) {
	t.TicketDescription = testcaseDescription
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

	// !!! Don't put any codes here. Preparation should be done in Testcase 0 !!!

	// TESTCASE 0
	// This testcase is used for preparation for this ticket.
	// If this testcase fails, the entire ticket will be skipped.
	tc := t.NewTestcase(0, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {
		// enter your testcase function here
		return tc.Failed()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)

	// TESTCASE 1
	tc = t.NewTestcase(1, "Enter your test case description here.")
	tc_func = func() models.TestcaseStatus {
		// enter your testcase function here
		return tc.Passed()
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

	// auto-register the new ticket in cmd/root.go init()
	if err := registerTicketInRoot(num); err != nil {
		fmt.Printf("Warning: could not auto-register ticket in root.go: %v\n", err)
	}
}

// registerTicketInRoot inserts the append line for the new ticket into root.go's init().
func registerTicketInRoot(num int) error {
	rootFile := filepath.Join("cmd", "root.go")
	data, err := os.ReadFile(rootFile)
	if err != nil {
		return err
	}

	insertLine := fmt.Sprintf("\tts = append(ts, new(tickets.Ticket%d))", num)
	content := string(data)

	// find the last occurrence of "ts = append(ts," and insert after that line
	lastAppend := lastIndexOf(content, "ts = append(ts,")
	if lastAppend < 0 {
		// fallback: insert after the comment anchor
		anchor := "// Add your tickets here"
		lastAppend = indexOf(content, anchor) + len(anchor)
	} else {
		// advance to end of that line
		for lastAppend < len(content) && content[lastAppend] != '\n' {
			lastAppend++
		}
	}

	updated := content[:lastAppend] + "\n" + insertLine + content[lastAppend:]
	if err := os.WriteFile(rootFile, []byte(updated), 0644); err != nil {
		return err
	}

	fmt.Printf("Registered Ticket%d in %s\n", num, rootFile)
	return nil
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func lastIndexOf(s, substr string) int {
	last := -1
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			last = i
		}
	}
	return last
}
