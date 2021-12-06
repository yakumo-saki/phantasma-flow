package logcollecter

import (
	"context"
)

type LogCollecterParamsBase struct {
	RunId  string
	JobId  string
	Alive  bool
	Ctx    context.Context
	Cancel context.CancelFunc
}
