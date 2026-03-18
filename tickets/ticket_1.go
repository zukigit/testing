package tickets

import "github.com/zukigit/testing/models"

type Ticket1 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []models.TestCase
}

func (t *Ticket1) NewTestcase(testcaseNo uint, testcaseDescription string) *models.TestCase {
	return models.NewTestcase(t.Ticket_no, testcaseNo, testcaseDescription)
}

func (t *Ticket1) GetTicketNo() uint {
	return t.Ticket_no
}

func (t *Ticket1) SetTicketNo(ticketNo uint) {
	t.Ticket_no = ticketNo
}

func (t *Ticket1) GetTicketDescription() string {
	return t.Ticket_description
}

func (t *Ticket1) SetTicketDescription(testcaseDescription string) {
	t.Ticket_description = testcaseDescription
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
