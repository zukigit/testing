package tickets

import "github.com/zukigit/testing/models"

type Ticket1 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []models.TestCase
}

func (t *Ticket1) New_testcase(testcase_id uint, testcase_description string) *models.TestCase {
	return models.New_testcase(testcase_id, testcase_description)
}

func (t *Ticket1) Get_ticket_no() uint {
	return t.Ticket_no
}

func (t *Ticket1) Set_ticket_no(ticket_no uint) {
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

func (t *Ticket1) Set_ticket_description(testcase_description string) {
	t.Ticket_description = testcase_description
}

func (t *Ticket1) Add_testcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket1) Get_testcases() []models.TestCase {
	return t.Testcases
}

func (t *Ticket1) Prepare() {
	t.Set_ticket_no(1)
	t.Set_ticket_description("Enter your ticket description here.")

	// TESTCASE 1
	tc := t.New_testcase(1, "Enter your test case description here.")
	tc_func := func() models.Testcase_status {
		//Enter your testcase function here
		return models.FAILED
	}
	tc.Set_function(tc_func)
	t.Add_testcase(tc)
}
