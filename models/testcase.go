package models

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type TestcaseStatus string
type LogAction string

const (
	PASSED     TestcaseStatus = "PASSED"
	FAILED     TestcaseStatus = "FAILED"
	MUST_CHECK TestcaseStatus = "MUST_CHECK"
	OUTPUT     LogAction      = "OUTPUT"
	STATUS     LogAction      = "STATUS"
)

type logEntry struct {
	Timestamp  string    `json:"timestamp"`
	Level      string    `json:"level"`
	TicketNo   uint      `json:"ticket_no,omitempty"`
	TestcaseNo uint      `json:"testcase_no,omitempty"`
	Action     LogAction `json:"action"`
	Message    string    `json:"message"`
}

type TestCase struct {
	Testcase_no          uint
	Testcase_description string
	Testcase_status      TestcaseStatus
	Duration             time.Duration
	function             func() TestcaseStatus
	Ticket_no            uint
	stdoutLogger         *log.Logger
	stderrLogger         *log.Logger
}

func NewTestcase(ticketNo, testcaseNo uint, testcase_description string) *TestCase {

	return &TestCase{
		Ticket_no:            ticketNo,
		Testcase_no:          testcaseNo,
		Testcase_description: testcase_description,
		stdoutLogger:         log.New(os.Stdout, "", 0),
		stderrLogger:         log.New(os.Stderr, "", 0),
	}
}

func (t *TestCase) GetTestcaseNo() uint {
	return t.Testcase_no
}

func (t *TestCase) GetTicketNo() uint {
	return t.Ticket_no
}

func (t *TestCase) SetTicketNo(ticket_no uint) {
	t.Ticket_no = ticket_no
}

func (t *TestCase) Get_ticket_description() string {
	return t.Testcase_description
}

func (t *TestCase) SetStatus(status TestcaseStatus) {
	t.Testcase_status = status
}

func (t *TestCase) Set_duration(duration time.Duration) {
	t.Duration = duration
}

func (t *TestCase) GetStatus() TestcaseStatus {
	return t.Testcase_status
}

func (t *TestCase) SetFunction(function func() TestcaseStatus) {
	t.function = function
}

func (t *TestCase) RunFunction() TestcaseStatus {
	return t.function()
}

func (t *TestCase) IsFunctionNil() bool {
	return t.function == nil
}

func (t *TestCase) writeLog(logger *log.Logger, level string, action LogAction, msg string) {
	entry := logEntry{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      level,
		Action:     action,
		TicketNo:   t.Ticket_no,
		TestcaseNo: t.Testcase_no,
		Message:    msg,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		logger.Printf(`{"level":"ERROR","msg":"failed to marshal log entry: %v"}`, err)
		return
	}
	logger.Println(string(data))
}

func (t *TestCase) InfoLog(msg string) {
	t.writeLog(t.stdoutLogger, "INFO", OUTPUT, msg)
}

func (t *TestCase) ErrorLog(msg string) {
	t.writeLog(t.stderrLogger, "ERROR", OUTPUT, msg)
}

func (t *TestCase) StatusLog(msg string) {
	t.writeLog(t.stdoutLogger, "INFO", STATUS, msg)
}

func (t *TestCase) Failed() TestcaseStatus {
	return FAILED
}

func (t *TestCase) Passed() TestcaseStatus {
	return PASSED
}

func (t *TestCase) MustCheck() TestcaseStatus {
	return MUST_CHECK
}
