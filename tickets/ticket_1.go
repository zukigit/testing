package tickets

import (
	"context"

	jaz_server "github.com/zukigit/testing/jaz_server"
	"github.com/zukigit/testing/lib"
	"github.com/zukigit/testing/models"
	"github.com/zukigit/testing/zabbix"
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

	// !!! Don't put any codes here. Preparation should be done in Testcase 0 !!!

	// TESTCASE 0
	// This testcase is used for preparation for this ticket.
	// If this testcase fails, the entire ticket will be skipped.
	tc0 := t.NewTestcase(0, "Enter your test case description here.")
	tc0.SetFunction(func() models.TestcaseStatus {
		envs := map[string]string{
			"ZABBIX_DB_TYPE":     string(models.DBTypePsql),
			"JAZ_SERVER_IMAGE":   string(models.JazServer609ImagePsql),
			"JAZ_DB_TYPE":        string(models.DBTypePsql),
			"JAZ_SERVER_VERSION": "1",
		}

		tc0.InfoLog("getting zabbix...")
		zbx, err := zabbix.NewZabbix(t.GetContext(), envs)
		if err != nil {
			tc0.ErrorLog("failed to get zabbix: %v", err)
			return tc0.Failed()
		}
		t.SetContext(context.WithValue(t.GetContext(), "zabbix-psql", zbx))

		tc0.InfoLog("getting jaz server...")
		_, err = jaz_server.NewJazServer(t.GetContext(), envs, zbx)
		if err != nil {
			tc0.ErrorLog("failed to get jaz server: %v", err)
			return tc0.Failed()
		}

		return tc0.Passed()
	})
	t.AddTestcase(tc0)

	// TESTCASE 1
	tc1 := t.NewTestcase(1, "Enter your test case description here.")
	tc1.SetFunction(func() models.TestcaseStatus {
		_, ok := lib.GetFromContext[zabbix.Zabbix](t.GetContext(), "zabbix-psql")
		if !ok {
			tc1.ErrorLog("failed to get zabbix from context")
			return tc1.Failed()
		}

		return tc1.Passed()
	})
	t.AddTestcase(tc1)
}
