package models

type Ticket interface {
	GetTicketNo() uint
	SetTicketNo(ticket_no uint)
	Get_ticket_description() string
	SetTicketDescription(testcase_description string)
	Set_PASSED_count(passed_count int)
	Set_FAILED_count(failed_count int)
	Set_MUSTCHECK_count(mustcheck_count int)
	AddTestcase(tc *TestCase)
	Prepare()
	Get_testcases() []TestCase
	NewTestcase(testcase_id uint, testcase_description string) *TestCase
}
