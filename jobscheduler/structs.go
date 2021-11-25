package jobscheduler

import (
	"container/list"
	"context"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
)

type JobScheduler struct {
	procman.ProcmanModuleStruct

	rootCtx    context.Context
	rootCancel context.CancelFunc
	jobs       map[string]job
	schedules  *list.List // list of schedule
	mutex      sync.Mutex
}

// Create from jobdefinition. Filter out not needed for scheduling.
type job struct {
	id      string
	name    string
	lastRun int64
	jobMeta objects.JobMetaInfo
}

type schedule struct {
	time  int64  // unixtime
	runId string // sha1 of uuid
	jobId string // job ID
}
