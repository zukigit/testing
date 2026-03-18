package tickets

import (
	"context"

	"github.com/zukigit/testing/models"
	"github.com/zukigit/testing/zabbix"
)

type Ticket2 struct {
	TicketNo          uint
	TicketDescription string
	Testcases         []models.TestCase
}

func (t *Ticket2) NewTestcase(testcaseNo uint, testcaseDescription string) *models.TestCase {
	return models.NewTestcase(t.TicketNo, testcaseNo, testcaseDescription)
}

func (t *Ticket2) GetTicketNo() uint {
	return t.TicketNo
}

func (t *Ticket2) SetTicketNo(ticketNo uint) {
	t.TicketNo = ticketNo
}

func (t *Ticket2) GetTicketDescription() string {
	return t.TicketDescription
}

func (t *Ticket2) SetTicketDescription(testcaseDescription string) {
	t.TicketDescription = testcaseDescription
}

func (t *Ticket2) AddTestcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket2) GetTestcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket2) Prepare() {
	t.SetTicketNo(2)
	t.SetTicketDescription("Enter your ticket description here.")

	// TESTCASE 1
	tc := t.NewTestcase(1, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {

		ctx := context.WithoutCancel(context.Background())
		defer ctx.Done()

		_, err := zabbix.NewZabbix(ctx)
		if err != nil {
			tc.ErrorLog("failed to get zabbix, err: %s", err.Error())
			return tc.Failed()
		}

		return tc.Passed()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)
}
