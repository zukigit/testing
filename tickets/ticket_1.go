package tickets

import "github.com/zukigit/testing/models"

type Ticket1 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []models.TestCase
}

func (t *Ticket1) NewTestcase(testcaseNo uint, testcase_description string) *models.TestCase {
	return models.NewTestcase(t.Ticket_no, testcaseNo, testcase_description)
}

func (t *Ticket1) GetTicketNo() uint {
	return t.Ticket_no
}

func (t *Ticket1) SetTicketNo(ticket_no uint) {
	t.Ticket_no = ticket_no
}

func (t *Ticket1) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket1) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket1) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket1) Get_ticket_description() string {
	return t.Ticket_description
}

func (t *Ticket1) SetTicketDescription(testcase_description string) {
	t.Ticket_description = testcase_description
}

func (t *Ticket1) AddTestcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket1) Get_testcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket1) Prepare() {
	t.SetTicketNo(1)
	t.SetTicketDescription("Enter your ticket description here.")

	// TESTCASE 1
	tc := t.NewTestcase(1, "Enter your test case description here.")
	tc_func := func() models.TestcaseStatus {
		//Enter your testcase function here
		return tc.Failed() // or tc.Passed() or tc.MustCheck()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)
}
