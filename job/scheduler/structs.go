package scheduler

import (
	"container/list"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/procman"
)

type JobScheduler struct {
	procman.ProcmanModuleStruct

	jobs      map[string]job
	runnables *list.List // list of schedule(runnable)
	schedules *list.List // list of schedule
	mutex     sync.Mutex
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

	scheduledAt int64
	queuedAt    int64
	runAt       int64
	endAt       int64
}
