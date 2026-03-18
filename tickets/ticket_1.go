package tickets

import (
	"context"

	"github.com/zukigit/testing/models"
)

type Ticket1 struct {
	TicketNo          uint
	TicketDescription string
	Testcases         []models.TestCase
	context           context.Context
}

func (t *Ticket1) SetContext(ctx context.Context) {
	t.context = ctx
}

func (t *Ticket1) GetContext() context.Context {
	return t.context
}

func (t *Ticket1) NewTestcase(testcaseNo uint, testcaseDescription string) *models.TestCase {
	return models.NewTestcase(t.TicketNo, testcaseNo, testcaseDescription)
}

func (t *Ticket1) GetTicketNo() uint {
	return t.TicketNo
}

func (t *Ticket1) SetTicketNo(ticketNo uint) {
	t.TicketNo = ticketNo
}

func (t *Ticket1) GetTicketDescription() string {
	return t.TicketDescription
}

func (t *Ticket1) SetTicketDescription(testcaseDescription string) {
	t.TicketDescription = testcaseDescription
}

func (t *Ticket1) AddTestcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket1) GetTestcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket1) Prepare() {
	t.SetTicketNo(1)
	t.SetTicketDescription("Enter your ticket description here.")
	t.SetContext(context.Background())

	// TESTCASE 1
	tc := t.NewTestcase(1, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {
		// enter your testcase function here
		return tc.Passed()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)
}
