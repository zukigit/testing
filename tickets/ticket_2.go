package tickets

import (
	"context"

	"github.com/zukigit/testing/models"
	"github.com/zukigit/testing/zabbix"
)

type Ticket2 struct {
	Ticket_no                                   uint
	Ticket_description                          string
	PASSED_count, FAILED_count, MUSTCHECK_count int
	Testcases                                   []models.TestCase
}

func (t *Ticket2) NewTestcase(testcaseNo uint, testcase_description string) *models.TestCase {
	return models.NewTestcase(t.Ticket_no, testcaseNo, testcase_description)
}

func (t *Ticket2) GetTicketNo() uint {
	return t.Ticket_no
}

func (t *Ticket2) SetTicketNo(ticket_no uint) {
	t.Ticket_no = ticket_no
}

func (t *Ticket2) Set_PASSED_count(passed_count int) {
	t.PASSED_count = passed_count
}

func (t *Ticket2) Set_FAILED_count(failed_count int) {
	t.FAILED_count = failed_count
}

func (t *Ticket2) Set_MUSTCHECK_count(mustcheck_count int) {
	t.MUSTCHECK_count = mustcheck_count
}

func (t *Ticket2) Get_ticket_description() string {
	return t.Ticket_description
}

func (t *Ticket2) SetTicketDescription(testcase_description string) {
	t.Ticket_description = testcase_description
}

func (t *Ticket2) AddTestcase(tc *models.TestCase) {
	t.Testcases = append(t.Testcases, *tc)
}

func (t *Ticket2) Get_testcases() []models.TestCase {
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

		zabbix, err := zabbix.NewZabbix(ctx)
		if err != nil {
			tc.ErrorLog("failed to get zabbix, err: %s", err.Error())
			return tc.Failed()
		}

		tc.InfoLog("DB Host: %s, DB Port: %s, DB Name: %s, DB Username: %s, DB Password: %s", zabbix.DBHost, zabbix.MappedPort, zabbix.DBName, zabbix.DBUsername, zabbix.DBPassword)
		return tc.Passed()
	}
	tc.SetFunction(tc_func)
	t.AddTestcase(tc)
}
