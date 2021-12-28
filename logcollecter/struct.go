package logcollecter

import (
	"context"
)

type LogCollecterParamsBase struct {
	RunId     string
	JobId     string
	JobNumber int
	Alive     bool
	Ctx       context.Context
	Cancel    context.CancelFunc
}
