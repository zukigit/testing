package models

type Ticket interface {
	Get_ticket_no() uint
	Set_ticket_no(ticket_no uint)
	Get_ticket_description() string
	Set_ticket_description(testcase_description string)
	Set_PASSED_count(passed_count int)
	Set_FAILED_count(failed_count int)
	Set_MUSTCHECK_count(mustcheck_count int)
	Add_testcase(tc *TestCase)
	Prepare()
	Get_testcases() []TestCase
	New_testcase(testcase_id uint, testcase_description string) *TestCase
}
