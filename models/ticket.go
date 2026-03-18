package models

import "context"

type Ticket interface {
	GetTicketNo() uint
	SetTicketNo(ticketNo uint)
	GetTicketDescription() string
	SetTicketDescription(TicketDescription string)
	AddTestcase(tc *TestCase)
	Prepare()
	GetTestcases() []TestCase
	NewTestcase(testcase_id uint, testcaseDescription string) *TestCase
	SetContext(ctx context.Context)
	GetContext() context.Context
}
