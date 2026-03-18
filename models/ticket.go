package models

type Ticket interface {
	GetTicketNo() uint
	SetTicketNo(ticketNo uint)
	GetTicketDescription() string
	SetTicketDescription(TicketDescription string)
	AddTestcase(tc *TestCase)
	Prepare()
	GetTestcases() []TestCase
	NewTestcase(testcase_id uint, testcaseDescription string) *TestCase
}
